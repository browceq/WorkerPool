package main

import (
	"context"
	"fmt"
	"sync"
)

// Интерфейсы, как вишенка на торте...

type WorkerPool interface {
	Work()
	AddWorker()
	RemoveWorker(id int)
	Wait()
}

func NewWorkerPool(ctx context.Context, numWorkers int, toProcessCh chan string) WorkerPool {

	pool := &workerPool{
		toProcessCh:    toProcessCh,
		ctx:            ctx,
		addWorkerCh:    make(chan struct{}),
		removeWorkerCh: make(chan int),
		workers:        make(map[int]Worker),
	}

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i)
		pool.workers[i] = worker
	}
	return pool

}

type workerPool struct {
	ctx context.Context
	wg  sync.WaitGroup

	toProcessCh    chan string
	addWorkerCh    chan struct{} // пустая структура вместо Worker, тк не надо передавать сам воркер
	removeWorkerCh chan int

	workers map[int]Worker // используется мапа для быстрого удаления воркеров

}

func (w *workerPool) Work() {

	//	Запускаем worker's при вызове функции
	for _, worker := range w.workers {
		w.wg.Add(1)
		go func(worker Worker) {
			defer w.wg.Done()
			worker.Work(w.ctx, w.toProcessCh)
		}(worker)
	}

	//	Запускаем select для обработки добавления или удаления воркеров
	go func() {
		for {

			select {
			case <-w.ctx.Done():
				return
			case <-w.addWorkerCh:
				w.addWorker()
			case id := <-w.removeWorkerCh:
				w.removeWorker(id)
			}

		}
	}()

}

func (w *workerPool) Wait() {
	w.wg.Wait()
}

// Внешние функции
func (w *workerPool) AddWorker() {
	w.addWorkerCh <- struct{}{}
}

func (w *workerPool) RemoveWorker(id int) {
	w.removeWorkerCh <- id
}

// Внутренние функции
func (w *workerPool) addWorker() {
	worker := NewWorker(len(w.workers))
	w.workers[len(w.workers)] = worker

	w.wg.Add(1)

	//	Запускаем созданный воркер
	go func() {
		defer w.wg.Done()
		worker.Work(w.ctx, w.toProcessCh)
	}()

}

func (w *workerPool) removeWorker(id int) {
	workerVal, ok := w.workers[id]
	if !ok {
		fmt.Printf("Failed to remove: there is no worker with id=%d\n", id)
		return
	}

	workerVal.Stop()
	delete(w.workers, id)

}
