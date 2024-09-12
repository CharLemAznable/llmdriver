package qwen

import (
	"context"
	"errors"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	loggerName = "qwen"

	defaultUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	configKeyForUrl          = "qwen.url"
	configKeyForApiKey       = "qwen.apiKey"
	configKeyFmtForApiKey    = "qwen.%s.apiKey"
	configKeyForWorkSpace    = "qwen.workSpace"
	configKeyFmtForWorkSpace = "qwen.%s.workSpace"

	headerAuthorization = "Authorization"
	headerWorkSpace     = "X-DashScope-WorkSpace"
	headerSSE           = "X-DashScope-SSE"

	headerSSEEnable = "enable"
)

var (
	logger = g.Log(loggerName)
	client = gx.Client().SetIntLog(logger).ContentJson()
)

type driver struct {
	name   string
	model  string
	client *gclientx.Client
	config *gmap.StrStrMap
}

func newDriver(name, model string) *driver {
	return &driver{
		name:   name,
		model:  model,
		client: client.Clone(),
		config: gmap.NewStrStrMap(true),
	}
}

func (d *driver) Available(ctx context.Context) bool {
	apiKey := d.getApiKey(ctx)
	if apiKey == "" {
		return false
	}
	authorization := "Bearer " + apiKey
	if workSpace := d.getWorkSpace(ctx); workSpace != "" {
		d.config.Sets(g.MapStrStr{
			headerAuthorization: authorization,
			headerWorkSpace:     workSpace,
		})
	} else {
		d.config.Sets(g.MapStrStr{
			headerAuthorization: authorization,
		})
	}
	return true
}

func (d *driver) getApiKey(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForApiKey, configKeyForApiKey)
}

func (d *driver) getWorkSpace(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForWorkSpace, configKeyForWorkSpace)
}

func (d *driver) Call(ctx context.Context, input llmdriver.Input) (output llmdriver.Output, err error) {
	req := buildReq(d.model, input, false)
	logger.Infof(ctx, "%s request once body: %s", d.model, gjson.MustEncodeString(req))
	rspContent, err := d.client.Header(d.config.Map()).PostContent(ctx, getUrl(ctx), req)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s response once body: %s", d.model, rspContent)
	return llmdriver.ParseJsonOutput(rspContent, getJsonKeysOptions()...)
}

func (d *driver) CallStream(ctx context.Context, input llmdriver.Input) (stream llmdriver.OutputStream) {
	outputStream := llmdriver.NewDefaultOutputStream()
	stream = outputStream

	req := buildReq(d.model, input, true)
	logger.Infof(ctx, "%s request stream body: %s", d.model, gjson.MustEncodeString(req))
	llmdriver.GoX(func() {
		eventSource := d.client.Header(d.config.Map()).
			Header(g.MapStrStr{headerSSE: headerSSEEnable}).
			PostEventSource(getUrl(ctx), req)
		defer func() {
			eventSource.Close()
			outputStream.Close(eventSource.Err())
		}()

		for event := range eventSource.Event() {
			if gstr.ToLower(event.Event) != "result" {
				outputStream.Close(errors.New(event.Data))
				return
			}
			logger.Infof(ctx, "%s response stream event: %s", d.model, gjson.MustEncodeString(event))
			outputEvent, err := llmdriver.ParseJsonOutputEvent(event, getJsonKeysOptions()...)
			if err != nil {
				outputStream.Close(err)
				return
			}
			outputStream.Push(outputEvent)
		}
	})
	return
}

func buildReq(model string, input llmdriver.Input, stream bool) g.Map {
	inputMap := g.Map{
		"messages": input.GetMessages(),
	}
	parameters := g.Map{
		"result_format": "message",
		"tools":         input.GetTools(),
	}
	if stream {
		parameters["incremental_output"] = true
	}
	return g.Map{
		"model":      model,
		"input":      inputMap,
		"parameters": parameters,
	}
}

func getUrl(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForUrl, defaultUrl)
}

func getJsonKeysOptions() []llmdriver.JsonKeysOption {
	return []llmdriver.JsonKeysOption{
		llmdriver.WithCodeKey("code"),
		llmdriver.WithMessageKey("message"),
		llmdriver.WithIdKey("request_id"),
		llmdriver.WithChoicesKey("output.choices"),
		llmdriver.WithChoicesMessageKey("message"),
		llmdriver.WithUsageKey("usage"),
		llmdriver.WithPromptKey("input_tokens"),
		llmdriver.WithCompletionKey("output_tokens"),
		llmdriver.WithTotalKey("total_tokens"),
	}
}
