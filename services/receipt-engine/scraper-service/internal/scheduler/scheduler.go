package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"backend_project/services/receipt-engine/scraper-service/internal/provider"
)

type Runner interface {
	Name() string
	Sync(ctx context.Context) error
}

type Task struct {
	Provider provider.Provider
	Interval time.Duration
	Type     provider.Type
}

type Scheduler struct {
	mu    sync.Mutex
	tasks []Task
	stop  chan struct{}
}

func New() *Scheduler {
	return &Scheduler{stop: make(chan struct{})}
}

func (s *Scheduler) Add(p provider.Provider, typ provider.Type) {
	s.mu.Lock()
	defer s.mu.Unlock()

	h := provider.Interval(typ)
	s.tasks = append(s.tasks, Task{
		Provider: p,
		Interval: time.Duration(h) * time.Hour,
		Type:     typ,
	})
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	tasks := make([]Task, len(s.tasks))
	copy(tasks, s.tasks)
	s.mu.Unlock()

	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			ticker := time.NewTicker(task.Interval)
			defer ticker.Stop()

			log.Printf("scheduler: started %s (interval=%v)", task.Provider.Name(), task.Interval)

			for {
				if err := task.Provider.Login(ctx, nil); err != nil {
					log.Printf("scheduler: %s login error: %v", task.Provider.Name(), err)
				}
				receipts, err := task.Provider.Sync(ctx)
				if err != nil {
					log.Printf("scheduler: %s sync error: %v", task.Provider.Name(), err)
				} else {
					log.Printf("scheduler: %s synced %d receipts", task.Provider.Name(), len(receipts))
				}

				select {
				case <-ticker.C:
				case <-ctx.Done():
					log.Printf("scheduler: %s stopped", task.Provider.Name())
					return
				}
			}
		}(t)
	}
	wg.Wait()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
