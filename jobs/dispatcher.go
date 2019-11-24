package jobs

import (
	"context"
	"log"
	"sync"
)

type Dispatcher struct {
	// A Queue of Job called by dispatcher
	JobQueue chan Job
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	// how many do work goroutine
	maxWorkers int
	ctx context.Context
	wg *sync.WaitGroup
}

func NewDispatcher(ctx context.Context, wg *sync.WaitGroup, maxWorkers int, jobQueue chan Job) *Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		WorkerPool: make(chan chan Job, maxWorkers),
		maxWorkers: maxWorkers,
		ctx:        ctx,
		wg:         wg,
	}
}

func (d *Dispatcher) Run() {
	d.wg.Add(d.maxWorkers+1)
	// starting n number of workers to do work
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.ctx, d.wg, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			log.Println("dispatch is Done.")
		case job := <-d.JobQueue:
			select {
			case <-d.ctx.Done():
				log.Println("blocking status but ctx.Done so dispatch is Done.")
				return
			case jobChannel := <-d.WorkerPool:
				jobChannel <- job
			}
		}
	}
}

