package zhipu

import "github.com/CharLemAznable/llmdriver"

func init() {
	modelMap := map[string]string{
		// 智谱清言
		"glm-4-plus":  "v4_plus",
		"glm-4-0520":  "v4_0520",
		"glm-4-long":  "v4_long",
		"glm-4-airx":  "v4_airx",
		"glm-4-air":   "v4_air",
		"glm-4-flash": "v4_flash",
	}
	for model, name := range modelMap {
		llmdriver.Register(model, newDriver(name, model))
	}
}
