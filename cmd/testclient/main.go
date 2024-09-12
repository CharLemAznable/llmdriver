package main

import (
	"context"
	"fmt"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/CharLemAznable/llmdriver/llmhttp"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	models = g.ArrayStr{
		"qwen-max",
		"moonshot-v1-8k",
		"glm-4-flash",
		"doubao-pro-4k",
		"hunyuan-functioncall",
	}
	client = gclientx.New().ContentJson().Prefix("http://127.0.0.1:38120")
	ctx    = context.Background()
)

func main() {
	for _, model := range models {
		promptWithModel(model)
		promptStreamWithModel(model)
	}
}

func promptWithModel(model string) {
	fmt.Printf("prompt with model: %s\n", model)
	content, _ := client.PostContent(ctx, "/completions", &llmhttp.Req{
		Model: model,
		Input: &llmhttp.Input{
			Prompt: llmdriver.String("介绍一下你自己"),
		},
	})
	rsp := new(llmhttp.Rsp)
	_ = gjson.New(content).Scan(rsp)
	if len(rsp.Output.Choices) > 0 {
		choice := rsp.Output.Choices[0]
		fmt.Println(
			llmdriver.StringValue(choice.Message.Content),
		)
	}
}

func promptStreamWithModel(model string) {
	fmt.Printf("prompt stream with model: %s\n", model)
	eventSource := client.PostEventSource("/completions", &llmhttp.Req{
		Model:  model,
		Stream: true,
		Input: &llmhttp.Input{
			Prompt: llmdriver.String("介绍一下你自己"),
		},
	})
	defer eventSource.Close()
	for event := range eventSource.Event() {
		output := new(llmhttp.Output)
		_ = gjson.New(event.Data).Scan(output)
		if len(output.Choices) > 0 {
			choice := output.Choices[0]
			fmt.Print(
				llmdriver.StringValue(choice.Message.Content),
			)
		}
	}
	fmt.Println()
}
