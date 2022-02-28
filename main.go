package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func worker(wg *sync.WaitGroup, jobs chan string) {
	defer wg.Done()

	for address := range jobs {
		uri, err := url.Parse(address)
		if err != nil {
			continue
		}
		if uri.Scheme == "" {
			uri.Scheme = "http"
		}

		resp, err := http.Get(uri.String())
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		hash := md5.Sum(body)
		fmt.Printf("%s %s\n", uri.String(), hex.EncodeToString(hash[:]))
	}
}

func processTasks(concurrencyLimit int, urls []string) {
	wg := &sync.WaitGroup{}
	inputs := make(chan string)

	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go worker(wg, inputs)
	}

	for _, uri := range urls {
		inputs <- uri
	}

	close(inputs)
	wg.Wait()
}

func main() {
	concurrencyLimit := flag.Int("parallel", 10, "max number of parallel HTTP requests. Min value is 1")
	flag.Parse()

	if *concurrencyLimit < 1 {
		*concurrencyLimit = 10
	}

	processTasks(*concurrencyLimit, flag.Args())
}
