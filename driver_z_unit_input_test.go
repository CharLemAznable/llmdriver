package llmdriver_test

import (
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/samber/lo"
	"testing"
)

func Test_Input_Messages(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		input := llmdriver.NewInput(make([]llmdriver.Message, 0), nil)
		t.AssertNE(input.GetMessages(), nil)
		t.Assert(len(input.GetMessages()), 0)
		t.AssertNil(input.GetTools())

		message := llmdriver.NewMessage(lo.ToPtr("role"), lo.ToPtr("content"))
		t.Assert(lo.FromPtr(message.GetRole()), "role")
		t.Assert(lo.FromPtr(message.GetContent()), "content")
		t.AssertNil(message.GetToolCallId())
		t.AssertNil(message.GetName())

		json := gjson.New(g.Slice{
			g.Map{
				"role":    "role",
				"content": "content",
				"tool_calls": g.Slice{
					g.Map{},
					g.Map{
						"type": "function",
						"function": g.Map{
							"name": "name",
						},
					},
				},
			},
		})
		messages, err := llmdriver.JsonToMessages(json)
		t.AssertNil(err)
		message = messages[0]
		t.Assert(lo.FromPtr(message.GetRole()), "role")
		t.Assert(lo.FromPtr(message.GetContent()), "content")
		toolCall := message.GetToolCalls()[0]
		t.AssertNil(toolCall.GetId())
		t.Assert(lo.FromPtr(toolCall.GetType()), "function")
		toolCallFunction := toolCall.GetFunction()
		t.Assert(lo.FromPtr(toolCallFunction.GetName()), "name")
		t.AssertNil(toolCallFunction.GetArguments())
		t.AssertNil(toolCall.GetIndex())
		t.AssertNil(message.GetToolCallId())
		t.AssertNil(message.GetName())
	})
}

func Test_Tools(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tool := llmdriver.NewTool(lo.ToPtr("function"),
			llmdriver.NewToolFunction(lo.ToPtr("name"),
				lo.ToPtr("description"),
				g.Map{"type": "string"}))
		t.Assert(lo.FromPtr(tool.GetType()), "function")
		t.Assert(lo.FromPtr(tool.GetFunction().GetName()), "name")
		t.Assert(lo.FromPtr(tool.GetFunction().GetDescription()), "description")
		t.Assert(tool.GetFunction().GetParameters(), g.Map{"type": "string"})

		json := gjson.New(g.Slice{
			g.Map{
				"type": "function",
				"function": g.Map{
					"name":        "name2",
					"description": "description2",
					"parameters":  g.Map{"type": "object"},
				},
			},
		})
		tools, err := llmdriver.JsonToTools(json)
		t.AssertNil(err)
		tool = tools[0]
		t.Assert(lo.FromPtr(tool.GetType()), "function")
		t.Assert(lo.FromPtr(tool.GetFunction().GetName()), "name2")
		t.Assert(lo.FromPtr(tool.GetFunction().GetDescription()), "description2")
		t.Assert(tool.GetFunction().GetParameters(), g.Map{"type": "object"})
	})
}
