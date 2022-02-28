package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func worker(wg *sync.WaitGroup, jobs chan string) {
	defer wg.Done()

	for address := range jobs {
		uri, err := url.Parse(address)
		if err != nil {
			log.Println(err)
		}
		if uri.Scheme == "" {
			uri.Scheme = "http"
		}

		resp, err := http.Get(uri.String())
		if err != nil {
			log.Printf("failed to make request to %s, error: %s\n", uri.String(), err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to process response body. URL: %s, error: %s\n", uri.String(), err)
			continue
		}

		hash := md5.Sum(body)
		log.Printf("%s %s\n", uri.String(), hex.EncodeToString(hash[:]))
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
	concurrencyLimit := flag.Int("parallel", 10, "")
	flag.Parse()

	processTasks(*concurrencyLimit, flag.Args())
}
