package main

import (
	_ "github.com/CharLemAznable/llmdriver/drivers/doubao"
	_ "github.com/CharLemAznable/llmdriver/drivers/hunyuan"
	_ "github.com/CharLemAznable/llmdriver/drivers/moonshot"
	_ "github.com/CharLemAznable/llmdriver/drivers/qwen"
	_ "github.com/CharLemAznable/llmdriver/drivers/zhipu"

	"github.com/CharLemAznable/llmdriver/llmhttp"
)

func main() {
	llmhttp.Server().Run()
}
