package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

// Meets the minimum specs
// Please see the README for more information

var crawled = make(map[string]bool)

//var crawledMutex sync.RWMutex
var final sync.WaitGroup
var host = ""

func main() {
	jobs := make(chan string)
	results := make(chan string)

	start := input()

	now := time.Now()
	go func() {
		jobs <- start
	}()

	var w worker = workerStruct{
		jobs:    jobs,
		results: results,
	}

	var wg sync.WaitGroup
	w.init(&wg)

	go func() {
		defer close(jobs)
		final.Wait()
	}()

	for r := range results {
		fmt.Println(r)
	}

	fmt.Println("All done in ", time.Since(now))
}

func input() string {
	var start string
	fmt.Print("Enter a starting URL (eg. https://nihalpandit.com) : ")
	_, err := fmt.Scanln(&start)
	if err != nil {
		panic(err)
	}

	URI, err := url.Parse(start)
	if err != nil || URI.Scheme == "mailto" || URI.Scheme == "tel" {
		panic(err)
	}

	host = URI.Host

	return start
}
