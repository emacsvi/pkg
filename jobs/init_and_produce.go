package jobs

import (
	"context"
	"github.com/emacsvi/lw-lotus/xgdata/common"
	"log"
	"sync"
	"time"
)

var (
	JobQueue chan Job
)

type PayLoad struct {
	block string
}

func (p *PayLoad) Do() error {
	time.Sleep(time.Second)
	log.Println("V: ", p.block)
	return nil
}

func LoopProduceData(ctx context.Context, wg *sync.WaitGroup, maxQueue, maxWorker int) {
	wg.Add(1)
	defer wg.Done()
	InitWorker(ctx, wg, maxQueue, maxWorker)
	for {
		msg := common.RandStringBytesMaskImprSrc(32)
		log.Println("P:", msg)
		job := &PayLoad{block: msg}

		select {
		case <-ctx.Done():
			log.Println("loopProdoceData ctx Done.")
			return
		case JobQueue <- job:
		}
	}
}

func InitWorker(ctx context.Context, wg *sync.WaitGroup, maxQueue, maxWorker int) {
	JobQueue = make(chan Job, maxQueue)
	d := NewDispatcher(ctx, wg, maxWorker, JobQueue)
	d.Run()
}
