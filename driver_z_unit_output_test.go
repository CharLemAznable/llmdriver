package llmdriver_test

import (
	"errors"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_Output_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		json := gjson.New(g.Map{
			"code":    "0",
			"message": "ok",
			"id":      "1",
			"choices": g.Slice{
				g.Map{
					"message": g.Map{
						"role":    "assistant",
						"content": "Hello!",
						"tool_calls": g.Slice{
							g.Map{
								"index": 1,
								"function": g.Map{
									"name":      "get_current_weather",
									"arguments": `{"location": "Boston, MA", "unit": "fahrenheit"}`,
								},
							},
						},
					},
				},
			},
			"usage": g.Map{
				"prompt":     1,
				"completion": 2,
			},
		})
		output, err := llmdriver.ParseJsonOutput(json.MustToJsonString(),
			llmdriver.WithCodeKey("code"),
			llmdriver.WithMessageKey("message"),
			llmdriver.WithIdKey("id"),
			llmdriver.WithChoicesKey("choices"),
			llmdriver.WithChoicesMessageKey("message"),
			llmdriver.WithUsageKey("usage"),
			llmdriver.WithPromptKey("prompt"),
			llmdriver.WithCompletionKey("completion"),
			llmdriver.WithTotalKey("total"))
		t.AssertNil(err)
		t.Assert(llmdriver.StringValue(output.Code()), "0")
		t.Assert(llmdriver.StringValue(output.Message()), "ok")
		t.Assert(llmdriver.StringValue(output.GetId()), "1")
		choice := output.GetChoices()[0]
		t.AssertNil(choice.GetIndex())
		t.Assert(llmdriver.IntValue(choice.GetIndex()), 0)
		t.AssertNil(choice.GetFinishReason())
		t.Assert(llmdriver.StringValue(choice.GetFinishReason()), "")
		t.Assert(llmdriver.StringValue(choice.GetMessage().GetRole()), "assistant")
		t.Assert(llmdriver.StringValue(choice.GetMessage().GetContent()), "Hello!")
		toolCall := choice.GetMessage().GetToolCalls()[0]
		t.AssertNil(toolCall.GetId())
		t.AssertNil(toolCall.GetType())
		t.AssertNil(choice.GetMessage().GetToolCallId())
		t.AssertNil(choice.GetMessage().GetName())
		t.Assert(llmdriver.IntValue(toolCall.GetIndex()), 1)
		t.Assert(llmdriver.StringValue(toolCall.GetFunction().GetName()), "get_current_weather")
		t.Assert(llmdriver.StringValue(toolCall.GetFunction().GetArguments()), `{"location": "Boston, MA", "unit": "fahrenheit"}`)
		usage := output.GetUsage()
		t.Assert(llmdriver.Int64Value(usage.GetPromptTokens()), 1)
		t.Assert(llmdriver.Int64Value(usage.GetCompletionTokens()), 2)
		t.AssertNil(usage.GetTotalTokens())
		t.Assert(llmdriver.Int64Value(usage.GetTotalTokens()), 0)

		output = llmdriver.NewJsonOutput(nil)
		t.AssertNil(output)

		json = gjson.New(g.Map{
			"id": "2",
			"choices": g.Slice{
				nil,
				g.Map{},
				g.Map{
					"delta": g.Map{
						"tool_calls": g.Slice{
							nil,
							g.Map{},
						},
					},
				},
			},
		})
		event := &gclientx.Event{
			Id:    "1",
			Event: "test",
			Data:  json.MustToJsonString(),
		}
		outputEvent, err := llmdriver.ParseJsonOutputEvent(event,
			llmdriver.WithEventId("0"),
			llmdriver.WithEvent(""),
			llmdriver.WithIdKey("id"),
			llmdriver.WithChoicesKey("choices"),
			llmdriver.WithChoicesMessageKey("delta"))
		t.AssertNil(err)
		t.Assert(llmdriver.StringValue(outputEvent.EventId()), "0")
		t.AssertNil(outputEvent.Event())
		t.Assert(llmdriver.StringValue(outputEvent.GetId()), "2")
		t.Assert(len(outputEvent.GetChoices()), 3)
		t.AssertNil(outputEvent.GetChoices()[0])
		choice = outputEvent.GetChoices()[1]
		t.AssertNil(choice.GetMessage())
		choice = outputEvent.GetChoices()[2]
		t.AssertNil(choice.GetMessage().GetToolCalls()[0])
		toolCall = choice.GetMessage().GetToolCalls()[1]
		t.AssertNil(toolCall.GetFunction())
		t.AssertNil(outputEvent.GetUsage())
		t.AssertNil(choice.GetMessage().GetToolCallId())
		t.AssertNil(choice.GetMessage().GetName())

		outputEvent = llmdriver.NewJsonOutputEvent(nil)
		t.AssertNil(outputEvent)
	})
}

func Test_Output_Stream(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		stream := llmdriver.NewDefaultOutputStream()
		go func() {
			stream.Push(llmdriver.NewJsonOutputEvent(gjson.New(g.Map{
				"id": "0",
			})))
			stream.Push(llmdriver.NewJsonOutputEvent(gjson.New(g.Map{
				"id": "1",
			})))
			stream.Close(errors.New("error"))
		}()

		event := <-stream.Event()
		t.Assert(llmdriver.StringValue(event.GetId()), "0")
		t.AssertNil(stream.Err())

		stream.Drain()
		t.Assert(stream.Err().Error(), "error")

		stream.Close(nil)
		t.Assert(stream.Err().Error(), "error")
	})
}
