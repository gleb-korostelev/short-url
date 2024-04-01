package worker

import (
	"context"
	"fmt"
	"sync"
)

type Task struct {
	Action func(ctx context.Context) error
	Done   chan struct{}
}

type DBWorkerPool struct {
	taskQueue  chan Task
	wg         sync.WaitGroup
	maxWorkers int
}

func NewDBWorkerPool(maxWorkers int) *DBWorkerPool {
	pool := &DBWorkerPool{
		taskQueue:  make(chan Task),
		maxWorkers: maxWorkers,
	}

	pool.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (p *DBWorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		if err := task.Action(context.Background()); err != nil {
			fmt.Printf("Error executing task: %v\n", err)
		}
		if task.Done != nil {
			close(task.Done)
		}
	}
}

func (p *DBWorkerPool) AddTask(task Task) {
	p.taskQueue <- task
}

func (p *DBWorkerPool) Shutdown() {
	close(p.taskQueue)
	p.wg.Wait()
}
