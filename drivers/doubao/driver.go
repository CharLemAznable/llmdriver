package doubao

import (
	"context"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmutex"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"io"
)

const (
	loggerName = "doubao"

	defaultBaseUrl = "https://ark.cn-beijing.volces.com/api/v3"
	defaultRegion  = "cn-beijing"

	configKeyForBaseUrl   = "doubao.url"
	configKeyForRegion    = "doubao.region"
	configKeyForApiKey    = "doubao.apiKey"
	configKeyFmtForApiKey = "doubao.%s.apiKey"
	configKeyForAK        = "doubao.ak"
	configKeyFmtForAK     = "doubao.%s.ak"
	configKeyForSK        = "doubao.sk"
	configKeyFmtForSK     = "doubao.%s.sk"
)

var (
	logger = g.Log(loggerName)
)

type driver struct {
	name     string
	model    string
	endpoint string
	mutex    *gmutex.RWMutex
	baseUrl  string
	region   string
	apiKey   string
	ak       string
	sk       string
	client   *arkruntime.Client
}

func newDriver(name, model, endpoint string) *driver {
	return &driver{
		name:     name,
		model:    model,
		endpoint: endpoint,
		mutex:    &gmutex.RWMutex{},
	}
}

func (d *driver) Available(ctx context.Context) bool {
	baseUrl, region := getBaseUrl(ctx), getRegion(ctx)
	apiKey, ak, sk := d.getApiKey(ctx), d.getAk(ctx), d.getSk(ctx)
	if apiKey != "" {
		d.updateClientWithApiKey(baseUrl, region, apiKey)
		return true
	} else if ak != "" && sk != "" {
		d.updateClientWithAkSk(baseUrl, region, ak, sk)
		return true
	}
	return false
}

func (d *driver) updateClientWithApiKey(baseUrl, region, apiKey string) {
	if baseUrl == d.baseUrl && region == d.region && apiKey == d.apiKey {
		return
	}
	d.mutex.LockFunc(func() {
		if baseUrl == d.baseUrl && region == d.region && apiKey == d.apiKey {
			return
		}
		d.baseUrl, d.region, d.apiKey = baseUrl, region, apiKey
		d.client = arkruntime.NewClientWithApiKey(apiKey,
			arkruntime.WithBaseUrl(baseUrl), arkruntime.WithRegion(region))
	})
}

func (d *driver) updateClientWithAkSk(baseUrl, region, ak, sk string) {
	if baseUrl == d.baseUrl && region == d.region && ak == d.ak && sk == d.sk {
		return
	}
	d.mutex.LockFunc(func() {
		if baseUrl == d.baseUrl && region == d.region && ak == d.ak && sk == d.sk {
			return
		}
		d.baseUrl, d.region, d.ak, d.sk = baseUrl, region, ak, sk
		d.client = arkruntime.NewClientWithAkSk(ak, sk,
			arkruntime.WithBaseUrl(baseUrl), arkruntime.WithRegion(region))
	})
}

func getBaseUrl(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForBaseUrl, defaultBaseUrl)
}

func getRegion(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForRegion, defaultRegion)
}

func (d *driver) getApiKey(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForApiKey, configKeyForApiKey)
}

func (d *driver) getAk(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForAK, configKeyForAK)
}

func (d *driver) getSk(ctx context.Context) (sk string) {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForSK, configKeyForSK)
}

func (d *driver) getClient() (client *arkruntime.Client) {
	d.mutex.RLockFunc(func() {
		client = d.client
	})
	return
}

func (d *driver) Call(ctx context.Context, input llmdriver.Input) (output llmdriver.Output, err error) {
	req := d.buildReq(input)
	logger.Infof(ctx, "%s request once body: %s", d.model, gjson.MustEncodeString(req))
	rsp, err := d.getClient().CreateChatCompletion(ctx, req)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s response once body: %s", d.model, gjson.MustEncodeString(rsp))
	return NewDoubaoOutput(rsp), nil
}

func (d *driver) CallStream(ctx context.Context, input llmdriver.Input) (stream llmdriver.OutputStream) {
	outputStream := llmdriver.NewDefaultOutputStream()
	stream = outputStream
	req := d.buildReq(input)
	logger.Infof(ctx, "%s request stream body: %s", d.model, gjson.MustEncodeString(req))
	llmdriver.GoX(func() {
		// optional set for return usage before [DONE] and after finish_reason:stop
		req.StreamOptions = &model.StreamOptions{IncludeUsage: true}
		stream, err := d.getClient().CreateChatCompletionStream(ctx, req)
		if err != nil {
			outputStream.Close(err)
			return
		}
		defer func() {
			_ = stream.Close()
			outputStream.Close(nil)
		}()
		for {
			streamRsp, err := stream.Recv()
			if err == io.EOF {
				outputStream.Close(nil)
				return
			}
			if err != nil {
				outputStream.Close(err)
				return
			}
			logger.Infof(ctx, "%s response stream event: %s", d.model, gjson.MustEncodeString(streamRsp))
			outputStream.Push(NewDoubaoOutputEvent(streamRsp))
		}
	})
	return
}

func (d *driver) buildReq(input llmdriver.Input) model.ChatCompletionRequest {
	request := model.ChatCompletionRequest{Model: d.endpoint}
	for _, message := range input.GetMessages() {
		request.Messages = append(request.Messages, &model.ChatCompletionMessage{
			Role: llmdriver.StringValue(message.GetRole()),
			Content: &model.ChatCompletionMessageContent{
				StringValue: message.GetContent(),
			},
			ToolCalls: gx.SliceMapping(message.GetToolCalls(), func(t llmdriver.ToolCall) *model.ToolCall {
				return &model.ToolCall{
					ID:   llmdriver.StringValue(t.GetId()),
					Type: model.ToolType(llmdriver.StringValue(t.GetType())),
					Function: model.FunctionCall{
						Name:      llmdriver.StringValue(t.GetFunction().GetName()),
						Arguments: llmdriver.StringValue(t.GetFunction().GetArguments()),
					},
				}
			}),
			ToolCallID: llmdriver.StringValue(message.GetToolCallId()),
		})
	}
	for _, tool := range input.GetTools() {
		request.Tools = append(request.Tools, &model.Tool{
			Type: model.ToolType(llmdriver.StringValue(tool.GetType())),
			Function: &model.FunctionDefinition{
				Name:        llmdriver.StringValue(tool.GetFunction().GetName()),
				Description: llmdriver.StringValue(tool.GetFunction().GetDescription()),
				Parameters:  tool.GetFunction().GetParameters(),
			},
		})
	}
	return request
}
