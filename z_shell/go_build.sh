#!/bin/sh

# llmdriver
echo ">>>>>>>> build package: llmdriver"
go build -v ./...

# llmdriver/drivers/doubao
echo ">>>>>>>> build package: llmdriver/drivers/doubao"
cd ./drivers/doubao
go build -v ./...
cd ./../..

# llmdriver/drivers/hunyuan
echo ">>>>>>>> build package: llmdriver/drivers/hunyuan"
cd ./drivers/hunyuan
go build -v ./...
cd ./../..

# llmdriver/drivers/moonshot
echo ">>>>>>>> build package: llmdriver/drivers/moonshot"
cd ./drivers/moonshot
go build -v ./...
cd ./../..

# llmdriver/drivers/qwen
echo ">>>>>>>> build package: llmdriver/drivers/qwen"
cd ./drivers/qwen
go build -v ./...
cd ./../..

# llmdriver/drivers/zhipu
echo ">>>>>>>> build package: llmdriver/drivers/zhipu"
cd ./drivers/zhipu
go build -v ./...
cd ./../..

# llmdriver/llmhttp
echo ">>>>>>>> build package: llmdriver/llmhttp"
cd ./llmhttp
go build -v ./...
cd ./..
