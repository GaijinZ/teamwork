package customerimporter

import (
	"sync"
	"sync/atomic"
)

type WorkerPool interface {
	Run()
	AddTask(task func())
	Wait()
	Close()
	Done() chan struct{}
}

type workerPool struct {
	sync.WaitGroup
	maxWorker   atomic.Int64
	queuedTaskC chan func()
	done        chan struct{}
}

func NewWorkerPool(maxWorker int64) WorkerPool {
	wp := workerPool{}
	wp.maxWorker.Store(maxWorker)
	wp.queuedTaskC = make(chan func(), 100)
	wp.done = make(chan struct{})
	return &wp
}

func (wp *workerPool) Run() {
	for i := 0; i < int(wp.maxWorker.Load()); i++ {
		go func(workerID int) {
			for task := range wp.queuedTaskC {
				task()
			}

			wp.maxWorker.Add(-1)
			if wp.maxWorker.Load() == 0 {
				wp.done <- struct{}{}
			}
		}(i + 1)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskC <- task
}

func (wp *workerPool) Close() {
	close(wp.queuedTaskC)
}

func (wp *workerPool) Done() chan struct{} {
	return wp.done
}
