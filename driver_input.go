package llmdriver

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gutil"
)

func NewInput(messages []Message, tools []Tool) Input {
	return &input{Messages: messages, Tools: tools}
}

type MessageOption func(*message)

func NewMessage(role, content *string, options ...MessageOption) Message {
	m := &message{
		Role:    role,
		Content: content,
	}
	for _, option := range options {
		option(m)
	}
	return m
}

func WithToolCalls(toolCalls []ToolCall) MessageOption {
	return func(m *message) {
		m.ToolCalls = toolCalls
	}
}

func WithToolCallId(toolCallId *string) MessageOption {
	return func(m *message) {
		m.ToolCallId = toolCallId
	}
}

func WithName(name *string) MessageOption {
	return func(m *message) {
		m.Name = name
	}
}

func NewTool(tType *string, function ToolFunction) Tool {
	return &tool{
		Type:     tType,
		Function: function,
	}
}

func NewToolFunction(name, description *string, parameters map[string]interface{}) ToolFunction {
	return &toolFunction{
		Name:        name,
		Description: description,
		Parameters:  gutil.MapCopy(parameters), // parameters参数可为空Map, 但不可为nil
	}
}

func NewToolCall(id, tType *string, function ToolCallFunction, index *int) ToolCall {
	return &toolCall{
		Id:       id,
		Type:     tType,
		Function: function,
		Index:    index,
	}
}

func NewToolCallFunction(name, arguments *string) ToolCallFunction {
	return &toolCallFunction{
		Name:      name,
		Arguments: arguments,
	}
}

func JsonToMessages(json *gjson.Json) (messages []Message, err error) {
	messagesArray := json.Array()
	for _, item := range messagesArray {
		messageJson := gjson.New(item)
		role := VarString(messageJson.Get("role"))
		content := VarString(messageJson.Get("content"))

		options := make([]MessageOption, 0)
		toolCalls := make([]ToolCall, 0)
		toolCallsArray := messageJson.Get("tool_calls").Array()
		for _, toolCallItem := range toolCallsArray {
			toolCallJson := gjson.New(toolCallItem)
			tType := toolCallJson.Get("type").String()
			if tType == "function" && !toolCallJson.Get("function").IsNil() {
				function := NewToolCallFunction(
					VarString(toolCallJson.Get("function.name")),
					VarString(toolCallJson.Get("function.arguments")))
				toolCalls = append(toolCalls, NewToolCall(
					VarString(toolCallJson.Get("id")), String(tType),
					function, VarInt(toolCallJson.Get("index"))))
			}
		}
		if len(toolCalls) > 0 {
			options = append(options, WithToolCalls(toolCalls))
		}
		options = append(options, WithToolCallId(VarString(messageJson.Get("tool_call_id"))))
		options = append(options, WithName(VarString(messageJson.Get("name"))))

		messages = append(messages, NewMessage(role, content, options...))
	}
	return
}

func JsonToTools(json *gjson.Json) (tools []Tool, err error) {
	toolsArray := json.Array()
	for _, item := range toolsArray {
		toolJson := gjson.New(item)
		tType := toolJson.Get("type").String()
		if tType == "function" && !toolJson.Get("function").IsNil() {
			function := NewToolFunction(
				VarString(toolJson.Get("function.name")),
				VarString(toolJson.Get("function.description")),
				toolJson.Get("function.parameters").Map())
			tools = append(tools, NewTool(String(tType), function))
		}
	}
	return
}

////////////////////////////////////////////////////////////////

type input struct {
	Messages []Message `json:"messages,omitempty"`
	Tools    []Tool    `json:"tools,omitempty"`
}

func (i *input) GetMessages() []Message {
	return i.Messages
}
func (i *input) GetTools() []Tool {
	return i.Tools
}

type message struct {
	Role       *string    `json:"role,omitempty"`
	Content    *string    `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallId *string    `json:"tool_call_id,omitempty"`
	Name       *string    `json:"name,omitempty"`
}

func (m *message) GetRole() *string {
	return m.Role
}
func (m *message) GetContent() *string {
	return m.Content
}
func (m *message) GetToolCalls() []ToolCall {
	return m.ToolCalls
}
func (m *message) GetToolCallId() *string {
	return m.ToolCallId
}
func (m *message) GetName() *string {
	return m.Name
}

type tool struct {
	Type     *string      `json:"type,omitempty"`
	Function ToolFunction `json:"function,omitempty"`
}

func (t *tool) GetType() *string {
	return t.Type
}
func (t *tool) GetFunction() ToolFunction {
	return t.Function
}

type toolFunction struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func (f *toolFunction) GetName() *string {
	return f.Name
}
func (f *toolFunction) GetDescription() *string {
	return f.Description
}
func (f *toolFunction) GetParameters() map[string]interface{} {
	return f.Parameters
}

type toolCall struct {
	Id       *string          `json:"id,omitempty"`
	Type     *string          `json:"type,omitempty"`
	Function ToolCallFunction `json:"function,omitempty"`
	Index    *int             `json:"index,omitempty"`
}

func (t *toolCall) GetId() *string {
	return t.Id
}
func (t *toolCall) GetType() *string {
	return t.Type
}
func (t *toolCall) GetFunction() ToolCallFunction {
	return t.Function
}
func (t *toolCall) GetIndex() *int {
	return t.Index
}

type toolCallFunction struct {
	Name      *string `json:"name,omitempty"`
	Arguments *string `json:"arguments,omitempty"`
}

func (f *toolCallFunction) GetName() *string {
	return f.Name
}
func (f *toolCallFunction) GetArguments() *string {
	return f.Arguments
}
