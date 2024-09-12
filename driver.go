package llmdriver

import "context"

type Driver interface {
	Available(ctx context.Context) (ok bool)
	Call(ctx context.Context, input Input) (output Output, err error)
	CallStream(ctx context.Context, input Input) (stream OutputStream)
}

type OutputStream interface {
	Event() <-chan OutputEvent
	Err() error
	Drain()
}

type Input interface {
	GetMessages() []Message
	GetTools() []Tool
}

type Output interface {
	Code() *string
	Message() *string
	GetId() *string
	GetChoices() []Choice
	GetUsage() Usage
}

type OutputEvent interface {
	EventId() *string
	Event() *string
	GetId() *string
	GetChoices() []Choice
	GetUsage() Usage
}

type Choice interface {
	GetIndex() *int
	GetFinishReason() *string
	GetMessage() Message
}

type Message interface {
	GetRole() *string
	GetContent() *string
	GetToolCalls() []ToolCall
	GetToolCallId() *string
	GetName() *string
}

type Tool interface {
	GetType() *string
	GetFunction() ToolFunction
}

type ToolFunction interface {
	GetName() *string
	GetDescription() *string
	GetParameters() map[string]interface{}
}

type ToolCall interface {
	GetId() *string
	GetType() *string
	GetFunction() ToolCallFunction
	GetIndex() *int
}

type ToolCallFunction interface {
	GetName() *string
	GetArguments() *string
}

type Usage interface {
	GetPromptTokens() *int64
	GetCompletionTokens() *int64
	GetTotalTokens() *int64
}
