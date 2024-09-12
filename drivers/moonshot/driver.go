package moonshot

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
	loggerName = "moonshot"

	defaultUrl = "https://api.moonshot.cn/v1/chat/completions"

	configKeyForUrl       = "moonshot.url"
	configKeyForApiKey    = "moonshot.apiKey"
	configKeyFmtForApiKey = "moonshot.%s.apiKey"

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
	req := buildReq(d.model, input)
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
	req := buildReq(d.model, input)
	req["stream"] = true
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
			outputEvent, err := llmdriver.ParseJsonOutputEvent(event,
				llmdriver.WithUsageKey("choices.0.usage")) // usage在choice中
			if err != nil {
				outputStream.Close(err)
				return
			}
			outputStream.Push(outputEvent)
		}
	})
	return
}

func buildReq(model string, input llmdriver.Input) g.Map {
	return g.Map{
		"model":    model,
		"messages": input.GetMessages(),
		"tools":    input.GetTools(),
	}
}

func getUrl(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForUrl, defaultUrl)
}
