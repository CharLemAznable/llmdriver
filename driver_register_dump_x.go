package llmdriver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/olekukonko/tablewriter"
)

func GetDriverMap(ctx context.Context) *gmap.TreeMap {
	driversMu.RLock()
	defer driversMu.RUnlock()
	m := gmap.NewTreeMap(gutil.ComparatorString, true)
	for k, v := range drivers {
		m.Set(k, v.Available(ctx))
	}
	return m
}

var (
	driverMapDump = gtype.NewBool()
)

func DumpDriverMap() {
	if !driverMapDump.Cas(false, true) {
		return
	}

	driverMap := GetDriverMap(context.Background())
	if driverMap.IsEmpty() {
		return
	}

	buffer := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{"MODEL", "AVAILABLE"})
	table.SetRowLine(false)
	table.SetBorder(true)
	table.SetCenterSeparator("|")

	driverMap.Iterator(func(key interface{}, value interface{}) bool {
		if gconv.Bool(value) {
			table.Append([]string{gconv.String(key), "ok"})
		} else {
			table.Append([]string{gconv.String(key), ""})
		}
		return true
	})
	table.Render()
	fmt.Printf("\n%s\n", buffer.String())
}
