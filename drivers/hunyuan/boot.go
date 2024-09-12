package hunyuan

import "github.com/CharLemAznable/llmdriver"

func init() {
	modelMap := map[string]string{
		// 腾讯混元
		"hunyuan-lite":          "lite",
		"hunyuan-standard":      "standard",
		"hunyuan-standard-256K": "standard_256K",
		"hunyuan-pro":           "pro",
		"hunyuan-code":          "code",
		"hunyuan-role":          "role",
		"hunyuan-functioncall":  "functioncall",
		"hunyuan-turbo":         "turbo",
	}
	for model, name := range modelMap {
		llmdriver.Register(model, newDriver(name, model))
	}
}
