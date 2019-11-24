package jobs

import (
	"context"
	"log"
	"sync"
)

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	ctx context.Context // this replace quit channel work
	wg *sync.WaitGroup // make sure all goroutine is quited
}

func NewWorker(ctx context.Context, wg *sync.WaitGroup, workerPool chan chan Job) *Worker {
	return &Worker{
		ctx: ctx,
		wg:wg,
		WorkerPool: workerPool,
		JobChannel: make(chan Job)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case <-w.ctx.Done():
				log.Println("work.ctx is Done")
				return
			default:
			}

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.Do(); err != nil {
					log.Fatalf("Error uploading to S3: %s", err.Error())
				}
			}
		}
	}()
}

