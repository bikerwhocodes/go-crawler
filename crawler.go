package main

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (worker workerStruct) crawl(weblink string) {
	resp, err := http.Get(weblink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return
	}
	processPage(resp.Body, worker)
	return
}

// NOTE: not my function, Modified from: https://vorozhko.net/get-all-links-from-html-page-with-go-lang
//Collect all links from response body and return it as an array of strings
func processPage(body io.Reader, worker workerStruct) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := normalize(attr.Val)
						uri, err := url.Parse(link)
						if err == nil && strings.Contains(uri.Host, host) && uri.Scheme != "mailto" && uri.Scheme != "tel" {
							if !crawled[link] {
								crawled[link] = true
								worker.enqueue(link)
								worker.result(link)
							}
						}
					}
				}
			}
		}
	}
}

func normalize(link string) string {
	suffix := "/"
	if strings.HasSuffix(link, suffix) {
		link = link[:len(link)-len(suffix)]
	}
	return link
}

//func read(index string) bool {
//	crawledMutex.RLock()
//	defer crawledMutex.RUnlock()
//	return crawled[index]
//}
//
//func write(index string) {
//	crawledMutex.Lock()
//	defer crawledMutex.Unlock()
//	crawled[index] = true
//}
