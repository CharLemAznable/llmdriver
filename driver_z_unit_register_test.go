package llmdriver_test

import (
	"context"
	"github.com/CharLemAznable/llmdriver"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_Driver_Register(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		err := g.Try(ctx, func(ctx context.Context) {
			llmdriver.Register("nil", nil)
		})
		t.Assert(err.Error(), "llmdriver: Register driver is nil")
		err = g.Try(ctx, func(ctx context.Context) {
			llmdriver.Register("unavailable", &unavailableDriver{})
		})
		t.AssertNil(err)
		err = g.Try(ctx, func(ctx context.Context) {
			llmdriver.Register("unavailable", &unavailableDriver{})
		})
		t.Assert(err.Error(), "llmdriver: Register called twice for driver unavailable")
		err = g.Try(ctx, func(ctx context.Context) {
			llmdriver.Register("available", &availableDriver{})
		})
		t.AssertNil(err)

		d, err := llmdriver.GetDriver(ctx, "nil")
		t.AssertNil(d)
		t.Assert(err.Error(), "llmdriver: unknown driver \"nil\" (forgotten import?)")
		d, err = llmdriver.GetDriver(ctx, "unavailable")
		t.AssertNil(d)
		t.Assert(err.Error(), "llmdriver: driver \"unavailable\" not available now")
		d, err = llmdriver.GetDriver(ctx, "available")
		t.AssertNE(d, nil)
		t.AssertNil(err)
	})
}

type unavailableDriver struct {
}

func (d *unavailableDriver) Available(_ context.Context) bool {
	return false
}

func (d *unavailableDriver) Call(_ context.Context, _ llmdriver.Input) (output llmdriver.Output, err error) {
	return
}

func (d *unavailableDriver) CallStream(_ context.Context, _ llmdriver.Input) (stream llmdriver.OutputStream) {
	return
}

type availableDriver struct {
}

func (d *availableDriver) Available(_ context.Context) bool {
	return true
}

func (d *availableDriver) Call(_ context.Context, _ llmdriver.Input) (output llmdriver.Output, err error) {
	return
}

func (d *availableDriver) CallStream(_ context.Context, _ llmdriver.Input) (stream llmdriver.OutputStream) {
	return
}
