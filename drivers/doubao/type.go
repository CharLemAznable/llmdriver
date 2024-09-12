package doubao

import (
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func NewDoubaoOutput(v model.ChatCompletionResponse) llmdriver.Output {
	return (doubaoOutput)(v)
}

func NewDoubaoOutputEvent(v model.ChatCompletionStreamResponse) llmdriver.OutputEvent {
	return (doubaoOutputEvent)(v)
}

func NewDoubaoChoice(v *model.ChatCompletionChoice) llmdriver.Choice {
	if v == nil {
		return nil
	}
	return (*doubaoChoice)(v)
}

func NewDoubaoMessage(v model.ChatCompletionMessage) llmdriver.Message {
	return (doubaoMessage)(v)
}

func NewDoubaoStreamChoice(v *model.ChatCompletionStreamChoice) llmdriver.Choice {
	if v == nil {
		return nil
	}
	return (*doubaoStreamChoice)(v)
}

func NewDoubaoStreamMessage(v model.ChatCompletionStreamChoiceDelta) llmdriver.Message {
	return (doubaoStreamMessage)(v)
}

func NewDoubaoToolCall(v *model.ToolCall) llmdriver.ToolCall {
	if v == nil {
		return nil
	}
	return (*doubaoToolCall)(v)
}

func NewDoubaoToolCallFunction(v model.FunctionCall) llmdriver.ToolCallFunction {
	return (doubaoToolCallFunction)(v)
}

func NewDoubaoUsage(v *model.Usage) llmdriver.Usage {
	if v == nil {
		return nil
	}
	return (*doubaoUsage)(v)
}

type doubaoOutput model.ChatCompletionResponse

func (v doubaoOutput) Code() *string {
	return nil
}
func (v doubaoOutput) Message() *string {
	return nil
}
func (v doubaoOutput) GetId() *string {
	return llmdriver.String(v.ID)
}
func (v doubaoOutput) GetChoices() []llmdriver.Choice {
	return gx.SliceMapping(v.Choices,
		func(t *model.ChatCompletionChoice) llmdriver.Choice { return NewDoubaoChoice(t) })
}
func (v doubaoOutput) GetUsage() llmdriver.Usage {
	return NewDoubaoUsage(&v.Usage)
}

type doubaoOutputEvent model.ChatCompletionStreamResponse

func (v doubaoOutputEvent) EventId() *string {
	return nil
}
func (v doubaoOutputEvent) Event() *string {
	return nil
}
func (v doubaoOutputEvent) GetId() *string {
	return llmdriver.String(v.ID)
}
func (v doubaoOutputEvent) GetChoices() []llmdriver.Choice {
	return gx.SliceMapping(v.Choices,
		func(t *model.ChatCompletionStreamChoice) llmdriver.Choice { return NewDoubaoStreamChoice(t) })
}
func (v doubaoOutputEvent) GetUsage() llmdriver.Usage {
	return NewDoubaoUsage(v.Usage)
}

type doubaoChoice model.ChatCompletionChoice

func (v *doubaoChoice) GetIndex() *int {
	return llmdriver.Int(v.Index)
}
func (v *doubaoChoice) GetFinishReason() *string {
	return llmdriver.StringNotEmpty(string(v.FinishReason))
}
func (v *doubaoChoice) GetMessage() llmdriver.Message {
	return NewDoubaoMessage(v.Message)
}

type doubaoMessage model.ChatCompletionMessage

func (v doubaoMessage) GetRole() *string {
	return llmdriver.String(v.Role)
}
func (v doubaoMessage) GetContent() *string {
	return v.Content.StringValue
}
func (v doubaoMessage) GetToolCalls() []llmdriver.ToolCall {
	return gx.SliceMapping(v.ToolCalls,
		func(t *model.ToolCall) llmdriver.ToolCall { return NewDoubaoToolCall(t) })
}
func (v doubaoMessage) GetToolCallId() *string {
	return llmdriver.String(v.ToolCallID)
}
func (v doubaoMessage) GetName() *string {
	return nil
}

type doubaoStreamChoice model.ChatCompletionStreamChoice

func (v *doubaoStreamChoice) GetIndex() *int {
	return llmdriver.Int(v.Index)
}
func (v *doubaoStreamChoice) GetFinishReason() *string {
	return llmdriver.StringNotEmpty(string(v.FinishReason))
}
func (v *doubaoStreamChoice) GetMessage() llmdriver.Message {
	return NewDoubaoStreamMessage(v.Delta)
}

type doubaoStreamMessage model.ChatCompletionStreamChoiceDelta

func (v doubaoStreamMessage) GetRole() *string {
	return llmdriver.String(v.Role)
}
func (v doubaoStreamMessage) GetContent() *string {
	return llmdriver.StringNotEmpty(v.Content)
}
func (v doubaoStreamMessage) GetToolCalls() []llmdriver.ToolCall {
	return gx.SliceMapping(v.ToolCalls,
		func(t *model.ToolCall) llmdriver.ToolCall { return NewDoubaoToolCall(t) })
}
func (v doubaoStreamMessage) GetToolCallId() *string {
	return nil
}
func (v doubaoStreamMessage) GetName() *string {
	return nil
}

type doubaoToolCall model.ToolCall

func (v *doubaoToolCall) GetId() *string {
	return llmdriver.String(v.ID)
}
func (v *doubaoToolCall) GetType() *string {
	return llmdriver.String(string(v.Type))
}
func (v *doubaoToolCall) GetFunction() llmdriver.ToolCallFunction {
	return NewDoubaoToolCallFunction(v.Function)
}
func (v *doubaoToolCall) GetIndex() *int {
	return nil
}

type doubaoToolCallFunction model.FunctionCall

func (v doubaoToolCallFunction) GetName() *string {
	return llmdriver.String(v.Name)
}
func (v doubaoToolCallFunction) GetArguments() *string {
	return llmdriver.String(v.Arguments)
}

type doubaoUsage model.Usage

func (v *doubaoUsage) GetPromptTokens() *int64 {
	return llmdriver.Int64(int64(v.PromptTokens))
}
func (v *doubaoUsage) GetCompletionTokens() *int64 {
	return llmdriver.Int64(int64(v.CompletionTokens))
}
func (v *doubaoUsage) GetTotalTokens() *int64 {
	return llmdriver.Int64(int64(v.TotalTokens))
}
