package status

import (
	"runtime"
	"time"
)

type Cpu struct {
	StatFile
}

func NewCpu() *Cpu {
	return &Cpu{*NewStatFile("/proc/stat")}
}

var (
	cpu_filter = func(s []string) []string { return s[1:4] }
	ncpu       = float64(runtime.NumCPU())
)

func (c *Cpu) Get(timeout time.Duration) Percent {

	c.Open()
	ti := time.Now()
	before := sum(c.ToInts(cpu_filter))
	c.Close()

	time.Sleep(timeout)

	c.Open()
	td := float64(time.Since(ti)) / float64(timeout)
	now := sum(c.ToInts(cpu_filter))
	c.Close()

	stat := float64(now - before)
	return Percent(((stat / (td * 100)) * 100) / ncpu)
}
