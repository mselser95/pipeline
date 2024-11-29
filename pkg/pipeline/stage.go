package pipeline

import (
	"context"
	"log"
	"sync"
)

// Stage represents a single stage in the pipeline.
type Stage struct {
	Name     string
	Inbound  chan any
	Outbound chan any
	Worker   func(ctx context.Context, item any) error
	ErrChan  chan error
}

// NewStage initializes a new Stage.
func NewStage(name string, worker func(ctx context.Context, item any) error, inbound, outbound chan any) *Stage {
	return &Stage{
		Name:     name,
		Inbound:  inbound,
		Outbound: outbound,
		Worker:   worker,
		ErrChan:  make(chan error, 1), // Buffered for error propagation
	}
}

// Run starts the stage's processing loop.
func (s *Stage) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] Context cancelled, shutting down", s.Name)
				return
			case item, ok := <-s.Inbound:
				if !ok {
					// Inbound channel closed, stop processing
					log.Printf("[%s] Inbound channel closed", s.Name)
					return
				}
				err := s.Worker(ctx, item)
				if err != nil {
					// Report error for centralized handling
					select {
					case s.ErrChan <- err:
					default:
						log.Printf("[%s] Error channel is full, dropping error: %v", s.Name, err)
					}
				}
				if s.Outbound != nil {
					select {
					case <-ctx.Done():
						log.Printf("[%s] Context cancelled while sending to outbound", s.Name)
						return
					case s.Outbound <- item:
					}
				}
			}
		}
	}()
}
