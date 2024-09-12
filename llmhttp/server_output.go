package llmhttp

import (
	"context"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func parseOutput(ctx context.Context, v llmdriver.Output) *Rsp {
	if g.IsNil(v) {
		return nil
	}
	return &Rsp{
		Code:    v.Code(),
		Message: v.Message(),
		Output: &Output{
			Id:      v.GetId(),
			Choices: gx.SliceMapping(v.GetChoices(), parseChoice),
			Usage:   parseUsage(v.GetUsage()),
			TraceId: gctx.CtxId(ctx),
		},
	}
}

func parseOutputEvent(ctx context.Context, v llmdriver.OutputEvent) *RspEvent {
	if g.IsNil(v) {
		return nil
	}
	return &RspEvent{
		EventId: v.EventId(),
		Event:   v.Event(),
		Output: &Output{
			Id:      v.GetId(),
			Choices: gx.SliceMapping(v.GetChoices(), parseChoice),
			Usage:   parseUsage(v.GetUsage()),
			TraceId: gctx.CtxId(ctx),
		},
	}
}

func parseChoice(v llmdriver.Choice) *Choice {
	if g.IsNil(v) {
		return nil
	}
	return &Choice{
		Index:        v.GetIndex(),
		FinishReason: v.GetFinishReason(),
		Message:      parseMessage(v.GetMessage()),
	}
}

func parseMessage(v llmdriver.Message) *RspMessage {
	if g.IsNil(v) {
		return nil
	}
	return &RspMessage{
		Message{
			Role:      v.GetRole(),
			Content:   v.GetContent(),
			ToolCalls: gx.SliceMapping(v.GetToolCalls(), parseToolCall),
		},
	}
}

func parseToolCall(v llmdriver.ToolCall) *ToolCall {
	if g.IsNil(v) {
		return nil
	}
	return &ToolCall{
		Id:       v.GetId(),
		Type:     v.GetType(),
		Function: parseToolCallFunction(v.GetFunction()),
		Index:    v.GetIndex(),
	}
}

func parseToolCallFunction(v llmdriver.ToolCallFunction) *ToolCallFunction {
	if g.IsNil(v) {
		return nil
	}
	return &ToolCallFunction{
		Name:      v.GetName(),
		Arguments: v.GetArguments(),
	}
}

func parseUsage(v llmdriver.Usage) *Usage {
	if g.IsNil(v) {
		return nil
	}
	return &Usage{
		PromptTokens:     v.GetPromptTokens(),
		CompletionTokens: v.GetCompletionTokens(),
		TotalTokens:      v.GetTotalTokens(),
	}
}
