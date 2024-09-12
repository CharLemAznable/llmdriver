package llmdriver

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("llmdriver: Register driver is nil")
	}
	if _, dup := drivers[gstr.ToLower(name)]; dup {
		panic("llmdriver: Register called twice for driver " + name)
	}
	drivers[gstr.ToLower(name)] = driver
}

func GetDriver(ctx context.Context, name string) (Driver, error) {
	driversMu.RLock()
	driver, ok := drivers[gstr.ToLower(name)]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("llmdriver: unknown driver %q (forgotten import?)", name)
	}
	ok = driver.Available(ctx)
	if !ok {
		return nil, fmt.Errorf("llmdriver: driver %q not available now", name)
	}
	return driver, nil
}
