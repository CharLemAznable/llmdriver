package llmdriver_test

import (
	"context"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_GoX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_ = g.Log()
		llmdriver.GoX(func() {
			panic("ignored")
		})
	})
}

func Test_GetConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		origin := g.Cfg().GetAdapter()
		defer g.Cfg().SetAdapter(origin)

		adapterContent, _ := gcfg.NewAdapterContent()
		g.Cfg().SetAdapter(adapterContent)

		v := llmdriver.GetConfigWithDefault(context.TODO(),
			"llm.driver.name", "openai")
		t.Assert(v, "openai")

		v = llmdriver.GetConfigWithNamePattern(context.TODO(),
			"test", "llm.driver.%s.name", "llm.driver.name")
		t.Assert(v, "")
	})
}
