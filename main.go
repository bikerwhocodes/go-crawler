package main

import (
	"fmt"
	"net/url"
	"time"
)

func main() {
	jobs := make(chan string, 32)
	results := make(chan string)

	start, host := input()
	jobs <- start

	fmt.Printf("Crawling %s\n", start)

	now := time.Now()

	var w worker = workerStruct{
		jobs:    jobs,
		results: results,
		baseurl: host,
		crawled: make(map[string]bool),
	}

	go w.init()

	for r := range results {
		fmt.Println(r)
	}

	fmt.Println("All done in ", time.Since(now))
}

func input() (string, string) {
	var start string
	fmt.Print("Enter a starting URL (eg. https://tredish.com) : ")
	_, err := fmt.Scanln(&start)
	if err != nil {
		panic(err)
	}

	URI, err := url.Parse(start)
	if err != nil || URI.Scheme == "mailto" || URI.Scheme == "tel" {
		panic(err)
	}

	return start, URI.Host
}
