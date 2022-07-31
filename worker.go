package main

import (
	"fmt"
	"runtime"
	"sync"
)

type worker interface {
	init(wg *sync.WaitGroup)
	enqueue(j string)
	result(result string)
	crawl(weblink string)
	workerFunc()
}

type workerStruct struct {
	jobs    chan string
	results chan string
}

func (worker workerStruct) init(wg *sync.WaitGroup) {
	fmt.Println("Initializing worker")
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.workerFunc()
		}()
	}

	go func() {
		defer close(worker.results)
		wg.Wait()
	}()
}

func (worker workerStruct) enqueue(job string) {
	worker.jobs <- job
}

func (worker workerStruct) result(result string) {
	worker.results <- result
}

func (worker workerStruct) workerFunc() {
	for j := range worker.jobs {
		final.Add(1)
		func(j string) {
			defer final.Done()
			worker.crawl(j)
		}(j)
	}
}
