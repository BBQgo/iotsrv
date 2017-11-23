.DEFAULT: build

.PHONY: build run

build:
	go build bbqsrv.go

run:
	go run bbqsrv.go
