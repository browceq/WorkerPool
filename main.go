package main

import (
	"context"
	"fmt"
	"time"
)

/*
Пример использования workerPool.
В примере использую контекст с таймаутом 20 секунд. 500 значений передается в WorkerPool
*/
func main() {

	toProcess := make(chan string, 500)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool := NewWorkerPool(ctx, 3, toProcess)

	go pool.Work()

	go func() {
		for i := 0; i < 500; i++ {
			toProcess <- fmt.Sprintf("Task %d", i)
		}
		close(toProcess)
	}()

	pool.AddWorker()
	pool.AddWorker()
	pool.RemoveWorker(1)

	pool.Wait()

}
