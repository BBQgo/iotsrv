.DEFAULT: build

.PHONY: build run vendor

vendor:
	glide update

build: vendor
	go build bbqsrv.go

run: 
	go run bbqsrv.go
