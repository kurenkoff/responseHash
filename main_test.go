package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_processTasks(t *testing.T) {
	rand.Seed(time.Now().Unix())

	type prepare func() ([]*httptest.Server, string)
	type args struct {
		concurrencyLimit int
		urls             []string
		prepare          prepare
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "no arguments",
			args: args{
				concurrencyLimit: 10,
				prepare: func() ([]*httptest.Server, string) {
					return []*httptest.Server{}, ""
				},
			},
		},
		{
			name: "single request",
			args: args{
				concurrencyLimit: 1,
				prepare: func() ([]*httptest.Server, string) {
					testBody := `{"test": true}`

					hash := md5.Sum([]byte(testBody))
					expectedHash := hex.EncodeToString(hash[:])

					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(testBody))
					}))

					return []*httptest.Server{ts}, fmt.Sprintf("%s %s\n", ts.URL, expectedHash)
				},
			},
		},
		{
			name: "invalid urls input",
			args: args{
				concurrencyLimit: 2,
				urls:             []string{"12345657", "test"},
				prepare: func() ([]*httptest.Server, string) {
					return []*httptest.Server{}, ""
				},
			},
		},
		{
			name: "5 requests with concurrency 2",
			args: args{
				concurrencyLimit: 2,
				prepare: func() ([]*httptest.Server, string) {
					var (
						numbersOfRequests = 5
						testBodyPattern   = `{"id": %s}`
						servers           = make([]*httptest.Server, numbersOfRequests)
						expectedBody      = ""
					)

					for i := 0; i < numbersOfRequests; i++ {
						body := fmt.Sprintf(testBodyPattern, rand.Int())
						servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(body))
						}))

						hash := md5.Sum([]byte(body))

						expectedBody += fmt.Sprintf(
							"%s %s\n",
							servers[i].URL,
							hex.EncodeToString(hash[:]),
						)
					}

					return servers, expectedBody
				},
			},
		},
		{
			name: "5 requests with concurrency 5",
			args: args{
				concurrencyLimit: 5,
				prepare: func() ([]*httptest.Server, string) {
					var (
						numbersOfRequests = 5
						testBodyPattern   = `{"reqNum": %s}`
						servers           = make([]*httptest.Server, numbersOfRequests)
						expectedBody      = ""
					)

					for i := 0; i < numbersOfRequests; i++ {
						body := fmt.Sprintf(testBodyPattern, rand.Int())
						servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(body))
						}))

						hash := md5.Sum([]byte(body))

						expectedBody += fmt.Sprintf(
							"%s %s\n",
							servers[i].URL,
							hex.EncodeToString(hash[:]),
						)
					}

					return servers, expectedBody
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServers, expected := tt.args.prepare()
			urls := make([]string, 0)
			if len(tt.args.urls) == 0 {
				for i := range testServers {
					urls = append(urls, testServers[i].URL)
				}
			}

			// not the best practice, but works in this case
			// do not want to make complex project structure for easy testing
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			processTasks(tt.args.concurrencyLimit, urls)

			w.Close()
			outRaw, _ := ioutil.ReadAll(r)
			os.Stdout = old

			expectedParts := strings.Split(strings.TrimSuffix(expected, "\n"), "\n")

			out := string(outRaw)

			for i := range expectedParts {
				if !strings.Contains(out, expectedParts[i]) {
					t.Errorf("Error in test case %s. Can't find %s in %s", tt.name, expectedParts[i], outRaw)
				}
			}

			outParts := strings.Split(strings.TrimSuffix(string(outRaw), "\n"), "\n")
			if len(outParts) > len(expectedParts) {
				t.Errorf(
					"Error in test case %s. Expected: %d results, got: %d",
					tt.name,
					len(expectedParts),
					len(outParts),
				)
			}

			for i := range testServers {
				testServers[i].Close()
			}

		})
	}
}

// Test_worker no scheme in URL case
func Test_worker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer ts.Close()

	expected := fmt.Sprintf("%s 098f6bcd4621d373cade4e832627b4f6\n", ts.URL)

	wg := &sync.WaitGroup{}
	input := make(chan string)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	wg.Add(1)
	go worker(wg, input)

	input <- strings.ReplaceAll(ts.URL, "http://", "")

	close(input)
	wg.Wait()

	w.Close()
	outRaw, _ := ioutil.ReadAll(r)
	os.Stdout = old

	if string(outRaw) != expected {
		t.Errorf("Error in no scheme test case. Expected: %s, got: %s", expected, string(outRaw))
	}
}
