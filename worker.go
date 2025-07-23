package main

import (
	"sync/atomic"
	"time"
)

type workerStruct struct {
	jobs    chan string
	results chan string
	baseurl string
	crawled map[string]bool
}

type worker interface {
	init()
	crawl(link string)
}

func (worker workerStruct) init() {
	var counter int32

	for {
		select {
		case link := <-worker.jobs:
			if worker.crawled[link] {
				continue
			}
			worker.crawled[link] = true

			atomic.AddInt32(&counter, 1)
			go func() {
				defer atomic.AddInt32(&counter, -1)
				worker.crawl(link)
			}()

			worker.results <- link

		case <-time.After(time.Millisecond * 10):
			if atomic.LoadInt32(&counter) == 0 {
				close(worker.results)
				return
			}
		}
	}
}
