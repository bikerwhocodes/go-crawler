package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/url"
	"time"
)

type lambdaInput struct {
	Start string `json:"start"`
}

type Response events.APIGatewayProxyResponse

func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	start := req.QueryStringParameters["start"]
	timeTaken, sites, err := crawler(start)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	body, err := json.Marshal(map[string]interface{}{
		"urlsCrawled": sites,
		"timeTaken":   timeTaken,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func crawler(start string) (string, []string, error) {
	jobs := make(chan string, 32)
	results := make(chan string)

	start, host := inputLambda(start)
	jobs <- start

	now := time.Now()

	var w worker = workerStruct{
		jobs:    jobs,
		results: results,
		baseurl: host,
		crawled: make(map[string]bool),
	}

	go w.init()

	sites := []string{}
	for r := range results {
		sites = append(sites, r)
	}

	return time.Since(now).String(), sites, nil
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
