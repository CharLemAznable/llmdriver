package llmhttp

import (
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/CharLemAznable/gfx/net/ghttpx"
	"github.com/CharLemAznable/gfx/net/gsse"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

func Server(name ...interface{}) (server *ghttpx.Server) {
	llmdriver.DumpDriverMap()
	server = gx.Server(name...)
	server.SetDumpRouterMap(false)
	server.BindHandler("POST:/completions", completions)
	server.SetDefaultAddr(":38120")
	return
}

func completions(request *ghttp.Request) {
	var (
		req = new(Req)
		err error
	)
	if err = request.ParseForm(req); err != nil {
		exitWithError(request, err)
	}
	driver, err := llmdriver.GetDriver(request.Context(), req.Model)
	if err != nil {
		exitWithError(request, err)
	}
	input, err := buildInput(request.Context(), req.Input)
	if err != nil {
		exitWithError(request, err)
	}
	if req.Stream {
		gsse.Handle(func(client *gsse.Client) {
			callStream(driver, input, client)
		})(request)
	} else {
		call(driver, input, request)
	}
}

func callStream(driver llmdriver.Driver, input llmdriver.Input, client *gsse.Client) {
	stream := driver.CallStream(client.Context(), input)
	defer stream.Drain()
	for outputEvent := range stream.Event() {
		rspEvent := parseOutputEvent(client.Context(), outputEvent)
		data := gjson.MustEncodeString(rspEvent.Output)
		client.SendEventWithId(
			llmdriver.StringValue(rspEvent.EventId),
			llmdriver.StringValue(rspEvent.Event), data)
	}
	if err := stream.Err(); err != nil {
		client.Response().WriteStatusExit(http.StatusBadRequest, err.Error())
	}
}

func call(driver llmdriver.Driver, input llmdriver.Input, request *ghttp.Request) {
	output, err := driver.Call(request.Context(), input)
	if err != nil {
		exitWithError(request, err)
	} else {
		request.Response.WriteJson(parseOutput(request.Context(), output))
	}
}

func exitWithError(request *ghttp.Request, err error) {
	request.Response.WriteStatusExit(http.StatusBadRequest, err.Error())
}
