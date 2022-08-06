package main

import (
	"net/http"
	"sort"
	"strings"
)

type result struct {
	query string
	index int
	res   http.Response
	err   error
}


func FetchContentParallel(urls []string, concurrencyLimit int) []result {
	semaphoreChan := make(chan struct{}, concurrencyLimit)
	resultsChan := make(chan *result)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, url := range urls {
		go func(i int, url string) {
			semaphoreChan <- struct{}{}

			res, err := http.Get(url)
			queryString := strings.Split(url, "=")
			r := &result{queryString[1], i, *res, err}

			resultsChan <- r
			<-semaphoreChan
		}(i, url)
	}

	var results []result
	for {
		r := <-resultsChan
		results = append(results, *r)

		if len(results) == len(urls) {
			break
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})
	return results
}
