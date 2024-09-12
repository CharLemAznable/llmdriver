package llmhttp

import (
	"context"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	invalidRequest = gerror.NewCode(gcode.CodeInvalidParameter, "invalid request")
)

func buildInput(ctx context.Context, input *Input) (llmdriver.Input, error) {
	if input == nil {
		return nil, invalidRequest
	}
	messages, err := parseMessages(ctx, input)
	if err != nil {
		return nil, err
	}
	tools, err := parseTools(ctx, input)
	if err != nil {
		return nil, err
	}
	if messages != nil {
		return llmdriver.NewInput(messages, tools), nil
	}
	return nil, invalidRequest
}

const (
	messageRoleUser = "user"
)

var (
	pMessageRoleUser = llmdriver.String(messageRoleUser)
)

func parseMessages(ctx context.Context, input *Input) ([]llmdriver.Message, error) {
	if input.MessagesTmpl != nil { // 使用模板表示会话列表
		messages, err := buildMessagesWithTmpl(ctx, input.MessagesTmpl)
		if err != nil {
			return nil, err
		}
		return messages, nil

	} else if len(input.Messages) > 0 { // 直传会话列表, 会话的Content字段可使用模板表示
		messages, err := buildMessages(ctx, input.Messages)
		if err != nil {
			return nil, err
		}
		return messages, nil

	} else if input.PromptTmpl != nil { // 使用模板表示提示词
		content, err := parseTmpl(ctx, input.PromptTmpl)
		if err != nil {
			return nil, err
		}
		return []llmdriver.Message{
			llmdriver.NewMessage(pMessageRoleUser, llmdriver.String(content)),
		}, nil

	} else if input.Prompt != nil { // 直传提示词
		return []llmdriver.Message{
			llmdriver.NewMessage(pMessageRoleUser, input.Prompt),
		}, nil
	}
	return nil, invalidRequest
}

func buildMessagesWithTmpl(ctx context.Context, tmpl *Tmpl) ([]llmdriver.Message, error) {
	content, err := parseTmpl(ctx, tmpl)
	if err != nil {
		return nil, err
	}
	messagesJson, err := gjson.LoadContent(content)
	if err != nil {
		return nil, err
	}
	return llmdriver.JsonToMessages(messagesJson)
}

func buildMessages(ctx context.Context, reqMessages []*ReqMessage) ([]llmdriver.Message, error) {
	var messages []llmdriver.Message
	for _, reqMessage := range reqMessages {
		if reqMessage == nil || reqMessage.Role == nil {
			continue
		}
		var content *string
		if reqMessage.ContentTmpl != nil {
			parsed, err := parseTmpl(ctx, reqMessage.ContentTmpl)
			if err != nil {
				return nil, err
			}
			content = llmdriver.String(parsed)
		} else if reqMessage.Content != nil {
			content = reqMessage.Content
		}

		options := make([]llmdriver.MessageOption, 0)
		toolCalls := make([]llmdriver.ToolCall, 0)
		for _, reqToolCall := range reqMessage.ToolCalls {
			if reqToolCall == nil || reqToolCall.Type == nil {
				continue
			}
			tType := llmdriver.StringValue(reqToolCall.Type)
			if tType == "function" && reqToolCall.Function != nil {
				toolCallFunction := llmdriver.NewToolCallFunction(
					reqToolCall.Function.Name,
					reqToolCall.Function.Arguments,
				)
				toolCall := llmdriver.NewToolCall(
					reqToolCall.Id,
					reqToolCall.Type,
					toolCallFunction,
					reqToolCall.Index,
				)
				toolCalls = append(toolCalls, toolCall)
			}
		}
		if len(toolCalls) > 0 {
			options = append(options, llmdriver.WithToolCalls(toolCalls))
		}
		options = append(options, llmdriver.WithToolCallId(reqMessage.ToolCallId))
		options = append(options, llmdriver.WithName(reqMessage.Name))

		messages = append(messages, llmdriver.NewMessage(reqMessage.Role, content, options...))
	}
	return messages, nil
}

func parseTools(ctx context.Context, input *Input) ([]llmdriver.Tool, error) {
	if input.ToolsTmpl != nil { // 使用模板表示工具列表
		tools, err := buildToolsWithTmpl(ctx, input.ToolsTmpl)
		if err != nil {
			return nil, err
		}
		return tools, nil
	} else if len(input.Tools) > 0 { // 直传工具列表
		return buildTools(input.Tools), nil
	}
	return nil, nil
}

func buildToolsWithTmpl(ctx context.Context, tmpl *Tmpl) ([]llmdriver.Tool, error) {
	content, err := parseTmpl(ctx, tmpl)
	if err != nil {
		return nil, err
	}
	toolsJson, err := gjson.LoadContent(content)
	if err != nil {
		return nil, err
	}
	return llmdriver.JsonToTools(toolsJson)
}

func buildTools(inputTools []*Tool) (tools []llmdriver.Tool) {
	for _, inputTool := range inputTools {
		if inputTool == nil || inputTool.Type == nil {
			continue
		}
		tType := llmdriver.StringValue(inputTool.Type)
		if tType == "function" && inputTool.Function != nil {
			tools = append(tools, llmdriver.NewTool(
				inputTool.Type,
				llmdriver.NewToolFunction(
					inputTool.Function.Name,
					inputTool.Function.Description,
					inputTool.Function.Parameters,
				),
			))
		}
	}
	return
}

func parseTmpl(ctx context.Context, tmpl *Tmpl) (string, error) {
	return gx.ViewX().Parse(ctx, tmpl.Name, tmpl.Params)
}
