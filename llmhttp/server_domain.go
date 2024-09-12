package llmhttp

type Message struct {
	Role       *string     `json:"role,omitempty"`
	Content    *string     `json:"content,omitempty"`
	ToolCalls  []*ToolCall `json:"tool_calls,omitempty"`
	ToolCallId *string     `json:"tool_call_id,omitempty"`
	Name       *string     `json:"name,omitempty"`
}

type ToolCall struct {
	Id       *string           `json:"id,omitempty"`
	Type     *string           `json:"type,omitempty"`
	Function *ToolCallFunction `json:"function,omitempty"`
	Index    *int              `json:"index,omitempty"`
}

type ToolCallFunction struct {
	Name      *string `json:"name,omitempty"`
	Arguments *string `json:"arguments,omitempty"`
}

////////////////////////////////////////////////////////////////

type Req struct {
	Model  string `v:"required"`
	Stream bool
	*Input
}

type Input struct {
	Messages     []*ReqMessage
	MessagesTmpl *Tmpl
	Prompt       *string
	PromptTmpl   *Tmpl

	Tools     []*Tool
	ToolsTmpl *Tmpl
}

type ReqMessage struct {
	Message
	ContentTmpl *Tmpl
}

type Tool struct {
	Type     *string
	Function *ToolFunction
}

type ToolFunction struct {
	Name        *string
	Description *string
	Parameters  map[string]interface{}
}

type Tmpl struct {
	Name   string
	Params map[string]interface{}
}

////////////////////////////////////////////////////////////////

type Rsp struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	*Output
}

type RspEvent struct {
	EventId *string
	Event   *string
	*Output
}

type Output struct {
	Id      *string   `json:"id,omitempty"`
	Choices []*Choice `json:"choices,omitempty"`
	Usage   *Usage    `json:"usage,omitempty"`
	TraceId string    `json:"trace_id,omitempty"`
}

type Choice struct {
	Index        *int        `json:"index,omitempty"`
	FinishReason *string     `json:"finish_reason,omitempty"`
	Message      *RspMessage `json:"message,omitempty"`
}

type RspMessage struct {
	Message
}

type Usage struct {
	PromptTokens     *int64 `json:"prompt_tokens,omitempty"`
	CompletionTokens *int64 `json:"completion_tokens,omitempty"`
	TotalTokens      *int64 `json:"total_tokens,omitempty"`
}
