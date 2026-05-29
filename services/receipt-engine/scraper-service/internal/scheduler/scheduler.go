package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
	"backend_project/services/receipt-engine/scraper-service/internal/provider"
)

type Task struct {
	Provider provider.Provider
	Interval time.Duration
	Type     provider.Type
}

type Scheduler struct {
	mu     sync.Mutex
	tasks  []Task
	stop   chan struct{}
	onSync func(ctx context.Context, providerName string, receipts []scrap.RawReceipt) error
}

func New() *Scheduler {
	return &Scheduler{stop: make(chan struct{})}
}

func (s *Scheduler) OnSync(f func(ctx context.Context, providerName string, receipts []scrap.RawReceipt) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onSync = f
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
				receipts, err := task.Provider.Sync(ctx)
				if err != nil {
					log.Printf("scheduler: %s sync error: %v", task.Provider.Name(), err)
				} else {
					log.Printf("scheduler: %s synced %d receipts", task.Provider.Name(), len(receipts))
					if s.onSync != nil {
						if err := s.onSync(ctx, task.Provider.Name(), receipts); err != nil {
							log.Printf("scheduler: %s onSync error: %v", task.Provider.Name(), err)
						}
					}
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
