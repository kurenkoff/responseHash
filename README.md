<h3 align="center">responseHash</h3>

## About The Project
Simple Go app which sends requests to specified URLs and prints MD5 hashes of response body's. 
No external dependencies required.

## Installation
Execute:
```shell
go install github.com/kurenkoff/responseHash
```

## Example
```shell
responseHash google.com twitter.com http://yandex.ru
```

## Arguments
```
-parallel - max number of parallel HTTP requests
```

## Tests
Execute:
```shell
go get github.com/kurenkoff/responseHash
cd $GOPATH/src/github.com/kurenkoff/responseHash
go test -race -cover ./...
```