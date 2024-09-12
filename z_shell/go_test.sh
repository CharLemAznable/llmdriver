#!/bin/sh

# llmdriver
echo ">>>>>>>> test package: llmdriver"
go test -v ./... -race -test.bench=.* -coverprofile=coverage.txt -covermode=atomic

# llmdriver/llmhttp
echo ">>>>>>>> test package: llmdriver/llmhttp"
cd ./llmhttp
go test -v ./... -race -test.bench=.* -coverprofile=coverage.txt -covermode=atomic
cd ./..
