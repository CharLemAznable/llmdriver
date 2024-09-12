package qwen

import "github.com/CharLemAznable/llmdriver"

func init() {
	modelMap := map[string]string{
		// 通义千问
		"qwen-turbo":           "turbo",
		"qwen-plus":            "plus",
		"qwen-max":             "max",
		"qwen-max-longcontext": "max_longcontext",

		// 通义千问-开源
		"qwen-1.8b-chat":             "open_1_8b",
		"qwen-1.8b-longcontext-chat": "open_1_8b_longcontext",
		"qwen-7b-chat":               "open_7b",
		"qwen-14b-chat":              "open_14b",
		"qwen-72b-chat":              "open_72b",

		// 通义千问1.5-开源
		"qwen1.5-0.5b-chat": "open1_5_0_5b",
		"qwen1.5-1.8b-chat": "open1_5_1_8b",
		"qwen1.5-7b-chat":   "open1_5_7b",
		"qwen1.5-14b-chat":  "open1_5_14b",
		"qwen1.5-32b-chat":  "open1_5_32b",
		"qwen1.5-72b-chat":  "open1_5_72b",
		"qwen1.5-110b-chat": "open1_5_110b",

		// 通义千问2-开源
		"qwen2-0.5b-instruct": "open2_0_5b",
		"qwen2-1.5b-instruct": "open2_1_5b",
		"qwen2-7b-instruct":   "open2_7b",
		"qwen2-72b-instruct":  "open2_72b",

		// 通义千问2-开源(57B规模14B激活参数的MOE模型)
		"qwen2-57b-a14b-instruct": "open2_57b_a14b",
	}
	for model, name := range modelMap {
		llmdriver.Register(model, newDriver(name, model))
	}
}
