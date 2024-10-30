package main

import (
	"context"
	"fmt"
)

// Интерфейсы, как вишенка на торте...

type Worker interface {
	Work(ctx context.Context, toProcessCh <-chan string)
	Stop()
}

type worker struct {
	id     int
	stopCh chan struct{}
}

func NewWorker(id int) Worker {
	return &worker{id, make(chan struct{})}
}

func (w *worker) Work(ctx context.Context, toProcessCh <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case val, ok := <-toProcessCh:
			if !ok {
				return
			}
			fmt.Printf("Worker: %d Value: %v\n", w.id, val)
		}
	}
}

func (w *worker) Stop() {
	w.stopCh <- struct{}{}
	close(w.stopCh)
}
