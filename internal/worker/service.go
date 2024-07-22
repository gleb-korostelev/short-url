// Package worker implements a worker pool that handles tasks concurrently.
// This package is designed to efficiently process tasks using multiple goroutines,
// managing synchronization and lifecycle of worker routines.
package worker

import (
	"context"
	"fmt"
	"sync"
)

// Task represents a unit of work to be executed by the worker pool.
// It contains an action to be executed and a channel to signal completion of the task.
type Task struct {
	Action func(ctx context.Context) error // Action is the function that performs the task.
	Done   chan struct{}                   // Done is used to signal the completion of the task.
}

// DBWorkerPool manages a pool of worker goroutines that execute Tasks.
type DBWorkerPool struct {
	taskQueue  chan Task      // taskQueue is a channel that holds tasks to be processed by the workers.
	wg         sync.WaitGroup // wg is used to wait for all workers to finish processing before shutdown.
	maxWorkers int            // maxWorkers defines the maximum number of worker goroutines.
}

// NewDBWorkerPool initializes a new DBWorkerPool with a specified number of workers.
// maxWorkers specifies the maximum number of concurrent workers in the pool.
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

// worker is a goroutine that processes Tasks from the taskQueue.
// It executes the Task's Action and signals completion via the Task's Done channel.
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

// AddTask submits a new Task to the pool. It adds the Task to the taskQueue.
func (p *DBWorkerPool) AddTask(task Task) {
	p.taskQueue <- task
}

// Shutdown gracefully stops the worker pool. It closes the taskQueue and waits for all workers to finish.
func (p *DBWorkerPool) Shutdown() {
	close(p.taskQueue)
	p.wg.Wait()
}
