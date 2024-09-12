package hunyuan

import (
	"context"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmutex"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

const (
	loggerName = "hunyuan"

	defaultScheme   = "HTTPS"
	defaultEndpoint = ""
	defaultRegion   = regions.Guangzhou

	configKeyForScheme   = "hunyuan.scheme"
	configKeyForEndpoint = "hunyuan.endpoint"
	configKeyForRegion   = "hunyuan.region"

	configKeyForSecretId     = "hunyuan.secretId"
	configKeyFmtForSecretId  = "hunyuan.%s.secretId"
	configKeyForSecretKey    = "hunyuan.secretKey"
	configKeyFmtForSecretKey = "hunyuan.%s.secretKey"
	configKeyForToken        = "hunyuan.token"
	configKeyFmtForToken     = "hunyuan.%s.token"
)

var (
	logger = g.Log(loggerName)
)

type driver struct {
	name  string
	model string

	mutex     *gmutex.RWMutex
	scheme    string
	endpoint  string
	region    string
	secretId  string
	secretKey string
	token     string
	client    *hunyuan.Client
}

func newDriver(name, model string) *driver {
	return &driver{
		name:  name,
		model: model,
		mutex: &gmutex.RWMutex{},
	}
}

func (d *driver) Available(_ context.Context) bool {
	scheme, endpoint, region := getScheme(context.Background()),
		getEndpoint(context.Background()), getRegion(context.Background())
	secretId, secretKey, token := d.getSecretId(context.Background()),
		d.getSecretKey(context.Background()), d.getToken(context.Background())
	if secretId != "" && secretKey != "" && token != "" {
		d.updateTokenClient(scheme, endpoint, region, secretId, secretKey, token)
		return true
	} else if secretId != "" && secretKey != "" {
		d.updateClient(scheme, endpoint, region, secretId, secretKey)
		return true
	}
	return false
}

func (d *driver) updateTokenClient(scheme, endpoint, region, secretId, secretKey, token string) {
	if scheme == d.scheme && endpoint == d.endpoint && region == d.region &&
		secretId == d.secretId && secretKey == d.secretKey && token == d.token {
		return
	}
	d.mutex.LockFunc(func() {
		if scheme == d.scheme && endpoint == d.endpoint && region == d.region &&
			secretId == d.secretId && secretKey == d.secretKey && token == d.token {
			return
		}
		d.scheme, d.endpoint, d.region, d.secretId, d.secretKey, d.token =
			scheme, endpoint, region, secretId, secretKey, token
		credential := tencent.NewTokenCredential(secretId, secretKey, token)
		clientProfile := profile.NewClientProfile()
		clientProfile.HttpProfile.Scheme = scheme
		clientProfile.HttpProfile.Endpoint = endpoint
		d.client, _ = hunyuan.NewClient(credential, region, clientProfile)
	})
}

func (d *driver) updateClient(scheme, endpoint, region, secretId, secretKey string) {
	if scheme == d.scheme && endpoint == d.endpoint && region == d.region &&
		secretId == d.secretId && secretKey == d.secretKey {
		return
	}
	d.mutex.LockFunc(func() {
		if scheme == d.scheme && endpoint == d.endpoint && region == d.region &&
			secretId == d.secretId && secretKey == d.secretKey {
			return
		}
		d.scheme, d.endpoint, d.region, d.secretId, d.secretKey =
			scheme, endpoint, region, secretId, secretKey
		credential := tencent.NewCredential(secretId, secretKey)
		clientProfile := profile.NewClientProfile()
		clientProfile.HttpProfile.Scheme = scheme
		clientProfile.HttpProfile.Endpoint = endpoint
		d.client, _ = hunyuan.NewClient(credential, region, clientProfile)
	})
}

func getScheme(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForScheme, defaultScheme)
}

func getEndpoint(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForEndpoint, defaultEndpoint)
}

func getRegion(ctx context.Context) string {
	return llmdriver.GetConfigWithDefault(ctx, configKeyForRegion, defaultRegion)
}

func (d *driver) getSecretId(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForSecretId, configKeyForSecretId)
}

func (d *driver) getSecretKey(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForSecretKey, configKeyForSecretKey)
}

func (d *driver) getToken(ctx context.Context) string {
	return llmdriver.GetConfigWithNamePattern(ctx, d.name, configKeyFmtForToken, configKeyForToken)
}

func (d *driver) getClient() (client *hunyuan.Client) {
	d.mutex.RLockFunc(func() {
		client = d.client
	})
	return
}

func (d *driver) Call(ctx context.Context, input llmdriver.Input) (output llmdriver.Output, err error) {
	req, err := d.buildReq(ctx, input)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s request once body: %s", d.model, req.ToJsonString())
	rsp, err := d.getClient().ChatCompletions(req)
	if err != nil {
		return
	}
	logger.Infof(ctx, "%s response once body: %s", d.model, rsp.ToJsonString())
	return NewHunyuanOutput(rsp), nil
}

func (d *driver) CallStream(ctx context.Context, input llmdriver.Input) (stream llmdriver.OutputStream) {
	outputStream := llmdriver.NewDefaultOutputStream()
	stream = outputStream
	req, err := d.buildReq(ctx, input)
	if err != nil {
		outputStream.Close(err)
		return
	}
	req.Stream = tencent.BoolPtr(true)
	logger.Infof(ctx, "%s request stream body: %s", d.model, req.ToJsonString())
	llmdriver.GoX(func() {
		rsp, err := d.getClient().ChatCompletions(req)
		if err != nil {
			outputStream.Close(err)
			return
		}
		defer func() {
			for range rsp.Events {
				// drain
			}
			outputStream.Close(nil)
		}()
		for event := range rsp.Events {
			if event.Err != nil {
				outputStream.Close(err)
				return
			}
			logger.Infof(ctx, "%s response stream event: %s", d.model, string(event.Data))
			data := &hunyuan.ChatCompletionsResponseParams{}
			if err = gjson.DecodeTo(event.Data, data); err != nil {
				outputStream.Close(err)
				return
			}
			outputStream.Push(NewHunyuanOutputEvent(data))
		}
	})
	return
}

func (d *driver) buildReq(ctx context.Context, input llmdriver.Input) (*hunyuan.ChatCompletionsRequest, error) {
	request := hunyuan.NewChatCompletionsRequest()
	request.SetContext(ctx)
	request.Model = llmdriver.String(d.model)
	for _, message := range input.GetMessages() {
		request.Messages = append(request.Messages, &hunyuan.Message{
			Role:       message.GetRole(),
			Content:    message.GetContent(),
			ToolCallId: message.GetToolCallId(),
			ToolCalls: gx.SliceMapping(message.GetToolCalls(), func(t llmdriver.ToolCall) *hunyuan.ToolCall {
				return &hunyuan.ToolCall{
					Id:   t.GetId(),
					Type: t.GetType(),
					Function: &hunyuan.ToolCallFunction{
						Name:      t.GetFunction().GetName(),
						Arguments: t.GetFunction().GetArguments(),
					},
				}
			}),
		})
	}
	for _, tool := range input.GetTools() {
		request.Tools = append(request.Tools, &hunyuan.Tool{
			Type: tool.GetType(),
			Function: &hunyuan.ToolFunction{
				Name:        tool.GetFunction().GetName(),
				Description: tool.GetFunction().GetDescription(),
				Parameters:  llmdriver.String(gjson.MustEncodeString(tool.GetFunction().GetParameters())),
			},
		})
	}
	return request, nil
}
