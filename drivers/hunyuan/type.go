package hunyuan

import (
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

func NewHunyuanOutput(v *hunyuan.ChatCompletionsResponse) llmdriver.Output {
	if v == nil {
		return nil
	}
	return (*hunyuanOutput)(v)
}

func NewHunyuanOutputEvent(v *hunyuan.ChatCompletionsResponseParams) llmdriver.OutputEvent {
	if v == nil {
		return nil
	}
	return (*hunyuanOutputEvent)(v)
}

func NewHunyuanChoice(v *hunyuan.Choice) llmdriver.Choice {
	if v == nil {
		return nil
	}
	return (*hunyuanChoice)(v)
}

func NewHunyuanMessage(v *hunyuan.Message) llmdriver.Message {
	if v == nil {
		return nil
	}
	return (*hunyuanMessage)(v)
}

func NewHunyuanStreamChoice(v *hunyuan.Choice) llmdriver.Choice {
	if v == nil {
		return nil
	}
	return (*hunyuanStreamChoice)(v)
}

func NewHunyuanStreamMessage(v *hunyuan.Delta) llmdriver.Message {
	if v == nil {
		return nil
	}
	return (*hunyuanStreamMessage)(v)
}

func NewHunyuanToolCall(v *hunyuan.ToolCall) llmdriver.ToolCall {
	if v == nil {
		return nil
	}
	return (*hunyuanToolCall)(v)
}

func NewHunyuanToolCallFunction(v *hunyuan.ToolCallFunction) llmdriver.ToolCallFunction {
	if v == nil {
		return nil
	}
	return (*hunyuanToolCallFunction)(v)
}

func NewHunyuanUsage(v *hunyuan.Usage) llmdriver.Usage {
	if v == nil {
		return nil
	}
	return (*hunyuanUsage)(v)
}

type hunyuanOutput hunyuan.ChatCompletionsResponse

func (v *hunyuanOutput) Code() *string {
	return nil
}
func (v *hunyuanOutput) Message() *string {
	return nil
}
func (v *hunyuanOutput) GetId() *string {
	return v.Response.Id
}
func (v *hunyuanOutput) GetChoices() []llmdriver.Choice {
	return gx.SliceMapping(v.Response.Choices,
		func(t *hunyuan.Choice) llmdriver.Choice { return NewHunyuanChoice(t) })
}
func (v *hunyuanOutput) GetUsage() llmdriver.Usage {
	return NewHunyuanUsage(v.Response.Usage)
}

type hunyuanOutputEvent hunyuan.ChatCompletionsResponseParams

func (v *hunyuanOutputEvent) EventId() *string {
	return nil
}
func (v *hunyuanOutputEvent) Event() *string {
	return nil
}
func (v *hunyuanOutputEvent) GetId() *string {
	return v.Id
}
func (v *hunyuanOutputEvent) GetChoices() []llmdriver.Choice {
	return gx.SliceMapping(v.Choices,
		func(t *hunyuan.Choice) llmdriver.Choice { return NewHunyuanStreamChoice(t) })
}
func (v *hunyuanOutputEvent) GetUsage() llmdriver.Usage {
	return NewHunyuanUsage(v.Usage)
}

type hunyuanChoice hunyuan.Choice

func (v *hunyuanChoice) GetIndex() *int {
	return nil
}
func (v *hunyuanChoice) GetFinishReason() *string {
	return v.FinishReason
}
func (v *hunyuanChoice) GetMessage() llmdriver.Message {
	return NewHunyuanMessage(v.Message)
}

type hunyuanMessage hunyuan.Message

func (v *hunyuanMessage) GetRole() *string {
	return v.Role
}
func (v *hunyuanMessage) GetContent() *string {
	return v.Content
}
func (v *hunyuanMessage) GetToolCalls() []llmdriver.ToolCall {
	return gx.SliceMapping(v.ToolCalls,
		func(t *hunyuan.ToolCall) llmdriver.ToolCall { return NewHunyuanToolCall(t) })
}
func (v *hunyuanMessage) GetToolCallId() *string {
	return v.ToolCallId
}
func (v *hunyuanMessage) GetName() *string {
	return nil
}

type hunyuanStreamChoice hunyuan.Choice

func (v *hunyuanStreamChoice) GetIndex() *int {
	return nil
}
func (v *hunyuanStreamChoice) GetFinishReason() *string {
	return v.FinishReason
}
func (v *hunyuanStreamChoice) GetMessage() llmdriver.Message {
	return NewHunyuanStreamMessage(v.Delta)
}

type hunyuanStreamMessage hunyuan.Delta

func (v *hunyuanStreamMessage) GetRole() *string {
	return v.Role
}
func (v *hunyuanStreamMessage) GetContent() *string {
	return v.Content
}
func (v *hunyuanStreamMessage) GetToolCalls() []llmdriver.ToolCall {
	return gx.SliceMapping(v.ToolCalls,
		func(t *hunyuan.ToolCall) llmdriver.ToolCall { return NewHunyuanToolCall(t) })
}
func (v *hunyuanStreamMessage) GetToolCallId() *string {
	return nil
}
func (v *hunyuanStreamMessage) GetName() *string {
	return nil
}

type hunyuanToolCall hunyuan.ToolCall

func (v *hunyuanToolCall) GetId() *string {
	return v.Id
}
func (v *hunyuanToolCall) GetType() *string {
	return v.Type
}
func (v *hunyuanToolCall) GetFunction() llmdriver.ToolCallFunction {
	return NewHunyuanToolCallFunction(v.Function)
}
func (v *hunyuanToolCall) GetIndex() *int {
	return nil
}

type hunyuanToolCallFunction hunyuan.ToolCallFunction

func (v *hunyuanToolCallFunction) GetName() *string {
	return v.Name
}
func (v *hunyuanToolCallFunction) GetArguments() *string {
	return v.Arguments
}

type hunyuanUsage hunyuan.Usage

func (v *hunyuanUsage) GetPromptTokens() *int64 {
	return v.PromptTokens
}
func (v *hunyuanUsage) GetCompletionTokens() *int64 {
	return v.CompletionTokens
}
func (v *hunyuanUsage) GetTotalTokens() *int64 {
	return v.TotalTokens
}
