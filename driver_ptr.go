package llmdriver

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/samber/lo"
)

func VarString(v *gvar.Var) *string {
	if v.IsNil() {
		return nil
	}
	return lo.ToPtr(v.String())
}

func VarInt(v *gvar.Var) *int {
	if v.IsNil() {
		return nil
	}
	return lo.ToPtr(v.Int())
}

func VarInt64(v *gvar.Var) *int64 {
	if v.IsNil() {
		return nil
	}
	return lo.ToPtr(v.Int64())
}
