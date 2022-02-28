<h3 align="center">responseHash</h3>

## About The Project
Simple Go app which sends requests to specified URLs and prints MD5 hashes of response body's. 
No external dependencies required.

Arguments of CLI tool are URLs. Invalid urls will be skipped. No logging. 

## Installation
Execute:
```shell
go install github.com/kurenkoff/responseHash@latest
```

## Example
```shell
responseHash google.com twitter.com http://yandex.ru

http://google.com f96740edb6e9ad8d32ebe0f7bb422951
http://yandex.ru 83e21774526d9038a26bc1484e9a3b2e
http://twitter.com 55a60e4dcbb5680fc1bc85f2ea0a74c4
```

## Arguments
```
-parallel - max number of parallel HTTP requests. Default value is 10. 
            Min value is 1 (if specified value is less than 1 app will use default value)
```

## Tests
Execute:
```shell
git clone https://github.com/kurenkoff/responseHash.git
cd responseHash
go test -race -cover ./.
```