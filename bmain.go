package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"net/url"
	"time"
)

type lambdaInput struct {
	Start string `json:"start"`
}

func main() {
	lambda.Start(crawler)
}

func crawler(event lambdaInput) (string, error) {
	jobs := make(chan string, 32)
	results := make(chan string)

	start, host := inputLambda(event.Start)
	jobs <- start

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

	return fmt.Sprintf("All done in %s", time.Since(now)), nil
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

func inputLambda(start string) (string, string) {
	URI, err := url.Parse(start)
	if err != nil || URI.Scheme == "mailto" || URI.Scheme == "tel" {
		panic(err)
	}

	return start, URI.Host
}
