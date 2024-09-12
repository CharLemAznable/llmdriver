package llmdriver

import (
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ParseJsonOutput(jsonString string, options ...JsonKeysOption) (Output, error) {
	json, err := gjson.DecodeToJson(jsonString)
	if err != nil {
		return nil, err
	}
	return NewJsonOutput(json, options...), nil
}

func ParseJsonOutputEvent(event *gclientx.Event, options ...JsonKeysOption) (OutputEvent, error) {
	json, err := gjson.DecodeToJson(event.Data)
	if err != nil {
		return nil, err
	}
	options = append([]JsonKeysOption{WithEventId(event.Id), WithEvent(event.Event)}, options...)
	return NewJsonOutputEvent(json, options...), nil
}

func NewJsonOutput(json *gjson.Json, options ...JsonKeysOption) Output {
	if json == nil {
		return nil
	}
	j := &jsonOutput{codeKey: "error.type", messageKey: "error.message",
		jsonKeys: &jsonKeys{Json: json, idKey: "id",
			choicesKey: "choices", choicesMessageKey: "message",
			usageKey: "usage", promptKey: "prompt_tokens",
			completionKey: "completion_tokens", totalKey: "total_tokens"}}
	for _, option := range options {
		option(j)
	}
	return j
}

func NewJsonOutputEvent(json *gjson.Json, options ...JsonKeysOption) OutputEvent {
	if json == nil {
		return nil
	}
	j := &jsonOutputEvent{id: nil, event: nil,
		jsonKeys: &jsonKeys{Json: json, idKey: "id",
			choicesKey: "choices", choicesMessageKey: "delta",
			usageKey: "usage", promptKey: "prompt_tokens",
			completionKey: "completion_tokens", totalKey: "total_tokens"}}
	for _, option := range options {
		option(j)
	}
	return j
}

func WithCodeKey(codeKey string) JsonKeysOption {
	return func(j jsonKeysOptional) {
		if output, ok := j.(*jsonOutput); ok {
			output.codeKey = codeKey
		}
	}
}

func WithMessageKey(messageKey string) JsonKeysOption {
	return func(j jsonKeysOptional) {
		if output, ok := j.(*jsonOutput); ok {
			output.messageKey = messageKey
		}
	}
}

func WithEventId(id string) JsonKeysOption {
	return func(j jsonKeysOptional) {
		if output, ok := j.(*jsonOutputEvent); ok {
			output.id = StringNotEmpty(id)
		}
	}
}

func WithEvent(event string) JsonKeysOption {
	return func(j jsonKeysOptional) {
		if output, ok := j.(*jsonOutputEvent); ok {
			output.event = StringNotEmpty(event)
		}
	}
}

func WithIdKey(idKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().idKey = idKey }
}

func WithChoicesKey(choicesKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().choicesKey = choicesKey }
}

func WithChoicesMessageKey(choicesMessageKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().choicesMessageKey = choicesMessageKey }
}

func WithUsageKey(usageKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().usageKey = usageKey }
}

func WithPromptKey(promptKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().promptKey = promptKey }
}

func WithCompletionKey(completionKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().completionKey = completionKey }
}

func WithTotalKey(totalKey string) JsonKeysOption {
	return func(j jsonKeysOptional) { j.GetJsonKeys().totalKey = totalKey }
}

func newJsonChoice(v interface{}, key ...string) Choice {
	if v == nil {
		return nil
	}
	messageKey := "message"
	if len(key) > 0 && len(key[0]) > 0 {
		messageKey = key[0]
	}
	return &jsonChoice{Json: gjson.New(v), messageKey: messageKey}
}

func newJsonMessage(v interface{}) Message {
	if v == nil {
		return nil
	}
	return &jsonMessage{Json: gjson.New(v)}
}

func newJsonToolCall(v interface{}) ToolCall {
	if v == nil {
		return nil
	}
	return &jsonToolCall{Json: gjson.New(v)}
}

func newJsonToolCallFunction(v interface{}) ToolCallFunction {
	if v == nil {
		return nil
	}
	return &jsonToolCallFunction{Json: gjson.New(v)}
}

func newJsonUsage(v interface{}, keys ...string) Usage {
	if v == nil {
		return nil
	}
	promptKey, completionKey, totalKey :=
		"prompt_tokens", "completion_tokens", "total_tokens"
	if len(keys) > 0 && len(keys[0]) > 0 {
		promptKey = keys[0]
	}
	if len(keys) > 1 && len(keys[1]) > 0 {
		completionKey = keys[1]
	}
	if len(keys) > 2 && len(keys[2]) > 0 {
		totalKey = keys[2]
	}
	return &jsonUsage{
		Json:          gjson.New(v),
		promptKey:     promptKey,
		completionKey: completionKey,
		totalKey:      totalKey,
	}
}

////////////////////////////////////////////////////////////////

type jsonKeysOptional interface {
	GetJsonKeys() *jsonKeys
}

type JsonKeysOption func(jsonKeysOptional)

type jsonOutput struct {
	codeKey    string
	messageKey string
	*jsonKeys
}

func (j *jsonOutput) Code() *string {
	return VarString(j.Get(j.codeKey))
}
func (j *jsonOutput) Message() *string {
	return VarString(j.Get(j.messageKey))
}
func (j *jsonOutput) GetJsonKeys() *jsonKeys {
	return j.jsonKeys
}

type jsonOutputEvent struct {
	id    *string
	event *string
	*jsonKeys
}

func (j *jsonOutputEvent) EventId() *string {
	return j.id
}
func (j *jsonOutputEvent) Event() *string {
	return j.event
}
func (j *jsonOutputEvent) GetJsonKeys() *jsonKeys {
	return j.jsonKeys
}

type jsonKeys struct {
	*gjson.Json
	idKey             string
	choicesKey        string
	choicesMessageKey string
	usageKey          string
	promptKey         string
	completionKey     string
	totalKey          string
}

func (j *jsonKeys) GetId() *string {
	return VarString(j.Get(j.idKey))
}
func (j *jsonKeys) GetChoices() []Choice {
	return gx.SliceMapping(j.Get(j.choicesKey).Array(),
		func(v interface{}) Choice { return newJsonChoice(v, j.choicesMessageKey) })
}
func (j *jsonKeys) GetUsage() Usage {
	return newJsonUsage(j.Get(j.usageKey).Val(), j.promptKey, j.completionKey, j.totalKey)
}

type jsonChoice struct {
	*gjson.Json
	messageKey string
}

func (j *jsonChoice) GetIndex() *int {
	return VarInt(j.Get("index"))
}
func (j *jsonChoice) GetFinishReason() *string {
	return VarString(j.Get("finish_reason"))
}
func (j *jsonChoice) GetMessage() Message {
	return newJsonMessage(j.Get(j.messageKey).Val())
}

type jsonMessage struct {
	*gjson.Json
}

func (j *jsonMessage) GetRole() *string {
	return VarString(j.Get("role"))
}
func (j *jsonMessage) GetContent() *string {
	return VarString(j.Get("content"))
}
func (j *jsonMessage) GetToolCalls() []ToolCall {
	return gx.SliceMapping(j.Get("tool_calls").Array(), newJsonToolCall)
}
func (j *jsonMessage) GetToolCallId() *string {
	return VarString(j.Get("tool_call_id"))
}
func (j *jsonMessage) GetName() *string {
	return VarString(j.Get("name"))
}

type jsonToolCall struct {
	*gjson.Json
}

func (j *jsonToolCall) GetId() *string {
	return VarString(j.Get("id"))
}
func (j *jsonToolCall) GetType() *string {
	return VarString(j.Get("type"))
}
func (j *jsonToolCall) GetFunction() ToolCallFunction {
	return newJsonToolCallFunction(j.Get("function").Val())
}
func (j *jsonToolCall) GetIndex() *int {
	return VarInt(j.Get("index"))
}

type jsonToolCallFunction struct {
	*gjson.Json
}

func (j *jsonToolCallFunction) GetName() *string {
	return VarString(j.Get("name"))
}
func (j *jsonToolCallFunction) GetArguments() *string {
	return VarString(j.Get("arguments"))
}

type jsonUsage struct {
	*gjson.Json
	promptKey     string
	completionKey string
	totalKey      string
}

func (j *jsonUsage) GetPromptTokens() *int64 {
	return VarInt64(j.Get(j.promptKey))
}
func (j *jsonUsage) GetCompletionTokens() *int64 {
	return VarInt64(j.Get(j.completionKey))
}
func (j *jsonUsage) GetTotalTokens() *int64 {
	return VarInt64(j.Get(j.totalKey))
}
