package doubao

import (
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	configKeyForEndpoints = "doubao.endpoints"
)

func init() {
	ctx := gctx.GetInitCtx()
	endpointMap := g.Cfg().MustGet(ctx, configKeyForEndpoints).MapStrVar()
	// 抖音火山豆包
	// 请求中的model字段仅支持填写推理接入点（形式为ep-xxxxxxxxx-yyyy）
	// 所以配置方式为:
	// 在endpoints配置中"自定义名称"=>{model:"模型名(可自定义)", endpoint:"推理接入点ID"}
	// 在其他配置中可以使用自定义名称单独配置客户端鉴权参数
	// 在调用时提交model: "模型名(可自定义)", 则使用对应的"推理接入点ID"作为实际请求中的model字段值
	for name, v := range endpointMap {
		modelEndpoint := v.MapStrStr()
		llmdriver.Register(modelEndpoint["model"],
			newDriver(name, modelEndpoint["model"], modelEndpoint["endpoint"]))
	}
}
