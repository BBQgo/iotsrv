.DEFAULT: build

.PHONY: build run vendor

vendor:
	glide update

build: vendor
	go build main.go

run:
	go run main.go
