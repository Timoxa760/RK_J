package x5club

import (
	"context"
	"sync"
)

type pageJob struct {
	page  int
	limit int
}

type pageResult struct {
	items []HistoryItem
	pages int
	err   error
}

type Pool struct {
	client  *Client
	workers int
}

func NewPool(client *Client, workers int) *Pool {
	if workers < 1 {
		workers = 1
	}
	return &Pool{client: client, workers: workers}
}

func (p *Pool) FetchAll(ctx context.Context, startPage, limit int) ([]HistoryItem, error) {
	firstItems, totalPages, err := p.client.GetHistory(startPage, limit)
	if err != nil {
		return nil, err
	}

	if totalPages <= 1 {
		return firstItems, nil
	}

	results := make([]HistoryItem, 0, totalPages*limit)
	results = append(results, firstItems...)

	jobs := make(chan pageJob, totalPages-1)
	resultsCh := make(chan pageResult, totalPages-1)
	var wg sync.WaitGroup

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				items, _, err := p.client.GetHistory(job.page, job.limit)
				resultsCh <- pageResult{items: items, err: err}
			}
		}()
	}

	for page := startPage + 1; page <= totalPages; page++ {
		jobs <- pageJob{page: page, limit: limit}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for res := range resultsCh {
		if res.err != nil {
			return nil, res.err
		}
		results = append(results, res.items...)
	}

	return results, nil
}
