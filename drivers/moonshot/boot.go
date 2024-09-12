package moonshot

import "github.com/CharLemAznable/llmdriver"

func init() {
	modelMap := map[string]string{
		// 月之暗面
		"moonshot-v1-8k":   "v1_8k",
		"moonshot-v1-32k":  "v1_32k",
		"moonshot-v1-128k": "v1_128k",
	}
	for model, name := range modelMap {
		llmdriver.Register(model, newDriver(name, model))
	}
}
