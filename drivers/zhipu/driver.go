package zhipu

import (
	"context"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/gfx/net/gclientx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	loggerName = "glm"

	defaultUrl = "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	configKeyForUrl       = "zhipu.url"
	configKeyForApiKey    = "zhipu.apiKey"
	configKeyFmtForApiKey = "zhipu.%s.apiKey"

	headerAuthorization = "Authorization"
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
	d.config.Sets(g.MapStrStr{
		headerAuthorization: authorization,
	})
	return true
}

func (d *driver) getApiKey(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForApiKey, configKeyForApiKey)
}

func (d *driver) Call(ctx context.Context, input llmdriver.Input) (output llmdriver.Output, err error) {
	req, err := d.buildReq(input, false)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s request once body: %s", d.model, gjson.MustEncodeString(req))
	rspContent, err := d.client.Header(d.config.Map()).PostContent(ctx, getUrl(ctx), req)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s response once body: %s", d.model, rspContent)
	return llmdriver.ParseJsonOutput(rspContent)
}

func (d *driver) CallStream(ctx context.Context, input llmdriver.Input) (stream llmdriver.OutputStream) {
	outputStream := llmdriver.NewDefaultOutputStream()
	stream = outputStream
	req, err := d.buildReq(input, true)
	if err != nil {
		outputStream.Close(err)
		return
	}
	logger.Infof(ctx, "%s request stream body: %s", d.model, gjson.MustEncodeString(req))
	llmdriver.GoX(func() {
		eventSource := d.client.Header(d.config.Map()).PostEventSource(getUrl(ctx), req)
		defer func() {
			eventSource.Close()
			outputStream.Close(eventSource.Err())
		}()

		for event := range eventSource.Event() {
			logger.Infof(ctx, "%s response stream event: %s", d.model, gjson.MustEncodeString(event))
			if event.Data == "[DONE]" {
				outputStream.Close(nil)
				return
			}
			outputEvent, err := llmdriver.ParseJsonOutputEvent(event)
			if err != nil {
				outputStream.Close(err)
				return
			}
			outputStream.Push(outputEvent)
		}
	})
	return
}

func (d *driver) buildReq(input llmdriver.Input, stream bool) (g.Map, error) {
	return g.Map{
		"model":    d.model,
		"messages": input.GetMessages(),
		"tools": gx.SliceMapping(input.GetTools(), func(t llmdriver.Tool) llmdriver.Tool {
			// glm模型的参数tools.function.parameters.type为必填参数
			if _, ok := t.GetFunction().GetParameters()["type"]; !ok {
				t.GetFunction().GetParameters()["type"] = "object"
			}
			// glm模型的参数tools.function.parameters.properties为必填参数
			if _, ok := t.GetFunction().GetParameters()["properties"]; !ok {
				t.GetFunction().GetParameters()["properties"] = g.Map{}
			}
			return t
		}),
		"stream": stream,
	}, nil
}

func getUrl(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForUrl, defaultUrl)
}
