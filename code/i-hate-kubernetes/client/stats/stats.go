package stats

import (
	"fmt"
	"runtime"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

type Stat struct {
	Memory             uint64
	MemoryPretty       string
	CpuUsage           *float64
	CpuUsagePercentage string
	Goroutines         int
}

func GetProcessStats(pid int, procStats *Stat) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	cpuUsage := GetCpu(pid)

	procStats.Goroutines = runtime.NumGoroutine()

	procStats.Memory = mem.Alloc
	procStats.MemoryPretty = console.PrettyMemoryAllocation(mem.Alloc)

	procStats.CpuUsage = cpuUsage
	procStats.CpuUsagePercentage = fmt.Sprintf("%.2f", *cpuUsage)
}
