#!/bin/sh

# llmdriver
echo ">>>>>>>> get package: llmdriver"
go get -t ./...

# llmdriver/drivers/doubao
echo ">>>>>>>> get package: llmdriver/drivers/doubao"
cd ./drivers/doubao
go get -t ./...
cd ./../..

# llmdriver/drivers/hunyuan
echo ">>>>>>>> get package: llmdriver/drivers/hunyuan"
cd ./drivers/hunyuan
go get -t ./...
cd ./../..

# llmdriver/drivers/moonshot
echo ">>>>>>>> get package: llmdriver/drivers/moonshot"
cd ./drivers/moonshot
go get -t ./...
cd ./../..

# llmdriver/drivers/qwen
echo ">>>>>>>> get package: llmdriver/drivers/qwen"
cd ./drivers/qwen
go get -t ./...
cd ./../..

# llmdriver/drivers/zhipu
echo ">>>>>>>> get package: llmdriver/drivers/zhipu"
cd ./drivers/zhipu
go get -t ./...
cd ./../..

# llmdriver/llmhttp
echo ">>>>>>>> get package: llmdriver/llmhttp"
cd ./llmhttp
go get -t ./...
cd ./..
