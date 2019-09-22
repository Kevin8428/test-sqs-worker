// Package workerpool
package workerpool

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Pool struct {
	Size             int
	WaitTime         int
	shutdownChannel  chan bool
	terminateChannel chan bool
	shuttingDown     bool
	running          int64
}

func (p *Pool) Start(worker Worker) {
	fmt.Println("003")
	workerChannel := make(chan Worker, p.Size)
	// Fill up the channel
	for i := 0; i < p.Size; i++ {
		workerChannel <- worker
	}
	fmt.Println("004")

	// No buffer. That means pool.Stop() will block until it stops.
	p.shutdownChannel = make(chan bool)
	p.terminateChannel = make(chan bool)
	fmt.Println("starting worker pool")
	for {
		select {
		case w := <-workerChannel:
			fmt.Println("005")
			if p.shuttingDown {
				break
			}
			// Keep track of how many are running. Use the sync.Mutex lock/unlock
			// to ensure atomicity of the "running" property
			go func() {
				fmt.Println("006")
				p.workerStarted()
				w.Work()
				workerChannel <- w
				p.workerFinished()
			}()
		case <-p.shutdownChannel:
			fmt.Println("shutting down channel")
			p.shuttingDown = true
		case <-p.terminateChannel:
			fmt.Println("terminating channel")
			return
		default:
			fmt.Println("007")
			time.Sleep(time.Duration(p.WaitTime) * time.Millisecond)
		}
	}
}

func (p *Pool) workerStarted() {
	fmt.Println("worker started")
}
func (p *Pool) workerFinished() {
	fmt.Println("worker finished")
}

// Stop working. Will block until the message is read by the working loop and # of
// running workers is down to zero
func (p *Pool) Shutdown() {
	p.shutdownChannel <- true
	for {
		running := atomic.LoadInt64(&p.running)

		if running == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
