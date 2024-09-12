package llmdriver

import (
	"context"
	"fmt"
	"github.com/CharLemAznable/gfx/container/gvarx"
	"github.com/CharLemAznable/gfx/frame/gx"
	"github.com/gogf/gf/v2/frame/g"
)

func LogError(err error) {
	g.Log().Errorf(context.Background(), "%+v", err)
}

func GoX(goroutineFunc func()) {
	gx.GoX(goroutineFunc, LogError)
}

func GetConfigWithDefault(ctx context.Context, pattern, def string) string {
	return gvarx.DefaultIfEmpty(g.Cfg().MustGet(ctx, pattern), def).String()
}

func GetConfigWithNamePattern(ctx context.Context, name, namePattern, defPattern string) (v string) {
	v = g.Cfg().MustGet(ctx, fmt.Sprintf(namePattern, name)).String()
	if v == "" {
		v = g.Cfg().MustGet(ctx, defPattern).String()
	}
	return
}
