package llmdriver

import "github.com/gogf/gf/v2/container/gvar"

func String(v string) *string {
	return &v
}

func StringNotEmpty(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func Int(v int) *int {
	return &v
}

func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func Int64(v int64) *int64 {
	return &v
}

func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

func VarString(v *gvar.Var) *string {
	if v.IsNil() {
		return nil
	}
	return String(v.String())
}

func VarInt(v *gvar.Var) *int {
	if v.IsNil() {
		return nil
	}
	return Int(v.Int())
}

func VarInt64(v *gvar.Var) *int64 {
	if v.IsNil() {
		return nil
	}
	return Int64(v.Int64())
}
