module github.com/CharLemAznable/llmdriver/cmd/server

go 1.20

require (
	github.com/CharLemAznable/llmdriver/drivers/doubao v0.1.0
	github.com/CharLemAznable/llmdriver/drivers/hunyuan v0.1.0
	github.com/CharLemAznable/llmdriver/drivers/moonshot v0.1.0
	github.com/CharLemAznable/llmdriver/drivers/qwen v0.1.0
	github.com/CharLemAznable/llmdriver/drivers/zhipu v0.1.0
	github.com/CharLemAznable/llmdriver/llmhttp v0.1.0
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/CharLemAznable/gfx v0.8.3 // indirect
	github.com/CharLemAznable/llmdriver v0.1.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogf/gf/v2 v2.7.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grokify/html-strip-tags-go v0.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1010 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan v1.0.1010 // indirect
	github.com/volcengine/volc-sdk-golang v1.0.23 // indirect
	github.com/volcengine/volcengine-go-sdk v1.0.158 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/CharLemAznable/llmdriver => ../..
	github.com/CharLemAznable/llmdriver/drivers/doubao => ../../drivers/doubao
	github.com/CharLemAznable/llmdriver/drivers/hunyuan => ../../drivers/hunyuan
	github.com/CharLemAznable/llmdriver/drivers/moonshot => ../../drivers/moonshot
	github.com/CharLemAznable/llmdriver/drivers/qwen => ../../drivers/qwen
	github.com/CharLemAznable/llmdriver/drivers/zhipu => ../../drivers/zhipu
	github.com/CharLemAznable/llmdriver/llmhttp => ../../llmhttp
)
