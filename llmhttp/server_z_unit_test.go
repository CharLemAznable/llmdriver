package llmhttp_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/CharLemAznable/gfx/os/gviewx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/CharLemAznable/llmdriver/llmhttp"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func init() {
	llmdriver.Register("echo", new(echoDriver))
	gx.ViewX().SetAdapter(gviewx.NewAdapterFile("z_test_data"))
}

type echoDriver struct {
}

func (d *echoDriver) Available(_ context.Context) bool {
	return true
}

func (d *echoDriver) Call(_ context.Context, input llmdriver.Input) (llmdriver.Output, error) {
	output := gjson.New(g.Map{"id": "echo"})
	var toolCalls g.Slice
	if len(input.GetTools()) > 0 {
		tool := input.GetTools()[0]
		toolCalls = g.Slice{
			g.Map{
				"type": tool.GetType(),
				"function": g.Map{
					"name": tool.GetFunction().GetName(),
				},
			},
		}
	}
	var choices g.List
	for _, message := range input.GetMessages() {
		if llmdriver.StringValue(message.GetContent()) == "error" {
			return nil, errors.New("error once")
		}
		choiceMessage := g.Map{
			"role":       message.GetRole(),
			"content":    message.GetContent(),
			"tool_calls": message.GetToolCalls(),
		}
		if toolCalls != nil {
			choiceMessage["tool_calls"] = toolCalls
		}
		choices = append(choices, g.Map{
			"message": choiceMessage,
		})
	}
	_ = output.Set("choices", choices)
	_ = output.Set("usage", g.Map{
		"prompt_tokens":     1,
		"completion_tokens": 1,
		"total_tokens":      2,
	})
	return llmdriver.NewJsonOutput(output), nil
}

func (d *echoDriver) CallStream(_ context.Context, input llmdriver.Input) llmdriver.OutputStream {
	outputStream := llmdriver.NewDefaultOutputStream()
	go func() {
		var toolCalls g.Slice
		if len(input.GetTools()) > 0 {
			tool := input.GetTools()[0]
			toolCalls = g.Slice{
				g.Map{
					"type": tool.GetType(),
					"function": g.Map{
						"name": tool.GetFunction().GetName(),
					},
				},
			}
		}
		for _, message := range input.GetMessages() {
			if llmdriver.StringValue(message.GetContent()) == "error" {
				outputStream.Close(errors.New("error stream"))
				return
			}
			choiceDelta := g.Map{
				"role":    message.GetRole(),
				"content": message.GetContent(),
			}
			if toolCalls != nil {
				choiceDelta["tool_calls"] = toolCalls
			}
			outputStream.Push(llmdriver.NewJsonOutputEvent(gjson.New(g.Map{
				"id": "echo",
				"choices": g.Slice{g.Map{
					"delta": choiceDelta,
				}},
			})))
		}
		outputStream.Close(nil)
	}()
	return outputStream
}

var (
	ctx    = context.TODO()
	urlFmt = `http://127.0.0.1:%d/completions`
	client = gclientx.New().ContentJson()
)

func Test_Req_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_error")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		resp, err := client.PostContent(ctx, url, g.Map{})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 The Model field is required")

		resp, err = client.PostContent(ctx, url, g.Map{"model": "nil"})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 llmdriver: unknown driver \"nil\" (forgotten import?)")

		resp, err = client.PostContent(ctx, url, g.Map{"model": "echo"})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 invalid request")

		resp, err = client.PostContent(ctx, url, g.Map{"model": "echo", "messages": g.Array{nil}})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 invalid request")

		resp, err = client.PostContent(ctx, url, g.Map{"model": "echo", "prompt": "error", "tools": g.Array{nil}})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 error once")

		eventSource := client.PostEventSource(url, g.Map{"model": "echo", "stream": "true", "prompt": "error"})
		defer eventSource.Close()
		event := <-eventSource.Event()
		t.AssertNil(event)
		t.Assert(eventSource.Err().Error(), "error stream")
	})
}

func Test_Req_Messages(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_messages")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		resp, err := client.PostContent(ctx, url, g.Map{
			"model": "echo",
			"messages": g.Array{
				g.Map{
					"role":    "assistant",
					"content": "hello",
					"tool_calls": g.List{
						g.Map{},
						g.Map{
							"type": "function",
							"function": g.Map{
								"name": "name",
							},
						},
					},
				},
				g.Map{
					"role": "user",
					"content_tmpl": g.Map{
						"name": "test_prompt",
						"params": g.Map{
							"Name": "John",
						},
					},
				},
			},
		})
		t.AssertNil(err)
		rsp := new(llmhttp.Rsp)
		_ = gjson.New(resp).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 2)
		t.Assert(rsp.Choices[0].Message.Role, "assistant")
		t.Assert(rsp.Choices[0].Message.Content, "hello")
		toolCall := rsp.Choices[0].Message.ToolCalls[0]
		t.Assert(toolCall.Type, "function")
		t.Assert(toolCall.Function.Name, "name")
		t.Assert(rsp.Choices[1].Message.Role, "user")
		t.Assert(rsp.Choices[1].Message.Content, "Hello, John!")
		t.Assert(rsp.Usage.PromptTokens, 1)
		t.Assert(rsp.Usage.CompletionTokens, 1)
		t.Assert(rsp.Usage.TotalTokens, 2)

		resp, err = client.PostContent(ctx, url, g.Map{
			"model": "echo",
			"messages": g.Array{
				g.Map{
					"role": "user",
					"content_tmpl": g.Map{
						"name": "error_prompt",
						"params": g.Map{
							"Name": "John",
						},
					},
				},
			},
		})
		t.Assert(resp, "")
		t.Assert(err.Error(), "400 template file \"error_prompt\" not found")
	})
}

func Test_Req_MessagesTmpl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_messages_tmpl")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		eventSource := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"messages_tmpl": g.Map{
				"name": "test_messages",
				"params": g.Map{
					"Name": "John",
				},
			},
		})
		defer eventSource.Close()
		event := <-eventSource.Event()
		t.AssertNE(event, nil)
		t.AssertNil(eventSource.Err())
		rsp := new(llmhttp.Output)
		_ = gjson.New(event.Data).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "assistant")
		t.Assert(rsp.Choices[0].Message.Content, "hello")
		event = <-eventSource.Event()
		t.AssertNE(event, nil)
		t.AssertNil(eventSource.Err())
		rsp = new(llmhttp.Output)
		_ = gjson.New(event.Data).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "user")
		t.Assert(rsp.Choices[0].Message.Content, "Hello, John!")

		eventSourceError := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"messages_tmpl": g.Map{
				"name": "error_messages",
				"params": g.Map{
					"Name": "John",
				},
			},
		})
		defer eventSourceError.Close()
		event = <-eventSourceError.Event()
		t.AssertNil(event)
		t.Assert(eventSourceError.Err().Error(), "template file \"error_messages\" not found")

		eventSourceIllegal := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"messages_tmpl": g.Map{
				"name": "illegal_messages",
				"params": g.Map{
					"Name": "John",
				},
			},
		})
		defer eventSourceIllegal.Close()
		event = <-eventSourceIllegal.Event()
		t.AssertNil(event)
		t.Assert(eventSourceIllegal.Err().Error(), "unsupported type \"\" for loading")
	})
}

func Test_Req_Prompt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_prompt")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		resp, err := client.PostContent(ctx, url, g.Map{
			"model":  "echo",
			"prompt": "hello",
		})
		t.AssertNil(err)
		rsp := new(llmhttp.Rsp)
		_ = gjson.New(resp).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "user")
		t.Assert(rsp.Choices[0].Message.Content, "hello")
		t.Assert(rsp.Usage.PromptTokens, 1)
		t.Assert(rsp.Usage.CompletionTokens, 1)
		t.Assert(rsp.Usage.TotalTokens, 2)
	})
}

func Test_Req_PromptTmpl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_prompt_tmpl")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		eventSource := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"prompt_tmpl": g.Map{
				"name": "test_prompt",
				"params": g.Map{
					"Name": "John",
				},
			},
		})
		defer eventSource.Close()
		event := <-eventSource.Event()
		t.AssertNE(event, nil)
		t.AssertNil(eventSource.Err())
		rsp := new(llmhttp.Output)
		_ = gjson.New(event.Data).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "user")
		t.Assert(rsp.Choices[0].Message.Content, "Hello, John!")

		eventSourceError := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"prompt_tmpl": g.Map{
				"name": "error_prompt",
				"params": g.Map{
					"Name": "John",
				},
			},
		})
		defer eventSourceError.Close()
		event = <-eventSourceError.Event()
		t.AssertNil(event)
		t.Assert(eventSourceError.Err().Error(), "template file \"error_prompt\" not found")
	})
}

func Test_Req_Tools(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_tools")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		resp, err := client.PostContent(ctx, url, g.Map{
			"model":  "echo",
			"prompt": "hello",
			"tools": g.Array{
				g.Map{
					"type": "function",
					"function": g.Map{
						"name":        "test_function",
						"description": "test function",
						"parameters": g.Map{
							"type": "object",
							"properties": g.Map{
								"Name": g.Map{
									"type": "string",
								},
							},
							"required": g.Array{"Name"},
						},
					},
				},
			},
		})
		t.AssertNil(err)
		rsp := new(llmhttp.Rsp)
		_ = gjson.New(resp).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "user")
		t.Assert(rsp.Choices[0].Message.Content, "hello")
		t.Assert(rsp.Choices[0].Message.ToolCalls[0].Type, "function")
		t.Assert(rsp.Choices[0].Message.ToolCalls[0].Function.Name, "test_function")
		t.Assert(rsp.Usage.PromptTokens, 1)
		t.Assert(rsp.Usage.CompletionTokens, 1)
		t.Assert(rsp.Usage.TotalTokens, 2)
	})
}

func Test_Req_ToolsTmpl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := llmhttp.Server("req_tools_tmpl")
		server.SetDumpRouterMap(false)
		_ = server.Start()
		defer func() { _ = server.Shutdown() }()
		url := fmt.Sprintf(urlFmt, server.GetListenedPort())

		eventSource := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"prompt": "hello",
			"tools_tmpl": g.Map{
				"name": "test_tools",
				"params": g.Map{
					"Key": "Name",
				},
			},
		})
		defer eventSource.Close()
		event := <-eventSource.Event()
		t.AssertNE(event, nil)
		t.AssertNil(eventSource.Err())
		rsp := new(llmhttp.Output)
		_ = gjson.New(event.Data).Scan(rsp)
		t.Assert(rsp.Id, "echo")
		t.Assert(len(rsp.Choices), 1)
		t.Assert(rsp.Choices[0].Message.Role, "user")
		t.Assert(rsp.Choices[0].Message.Content, "hello")
		t.Assert(rsp.Choices[0].Message.ToolCalls[0].Type, "function")
		t.Assert(rsp.Choices[0].Message.ToolCalls[0].Function.Name, "test_function")

		eventSourceError := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"prompt": "hello",
			"tools_tmpl": g.Map{
				"name": "error_tools",
				"params": g.Map{
					"Key": "Name",
				},
			},
		})
		defer eventSourceError.Close()
		event = <-eventSourceError.Event()
		t.AssertNil(event)
		t.Assert(eventSourceError.Err().Error(), "template file \"error_tools\" not found")

		eventSourceIllegal := client.PostEventSource(url, g.Map{
			"model":  "echo",
			"stream": "true",
			"prompt": "hello",
			"tools_tmpl": g.Map{
				"name": "illegal_tools",
				"params": g.Map{
					"Key": "Name",
				},
			},
		})
		defer eventSourceIllegal.Close()
		event = <-eventSourceIllegal.Event()
		t.AssertNil(event)
		t.Assert(eventSourceIllegal.Err().Error(), "unsupported type \"\" for loading")
	})
}
