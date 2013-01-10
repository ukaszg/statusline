package status

import (
	"runtime"
	"time"
)

var (
	cpu_filter = func(s []string) []string { return s[1:4] }
	ncpu       = float64(runtime.NumCPU())
)

type Cpu struct {
	StatFile
	Timeout time.Duration
}

func NewCpu(timeout time.Duration) *Cpu {
	return &Cpu{*NewStatFile("/proc/stat"), timeout}
}
func (c *Cpu) Get() Percent {
	c.Open()
	ti := time.Now()
	before := sum(c.ToInts(cpu_filter))
	c.Close()

	time.Sleep(c.Timeout)

	c.Open()
	td := float64(time.Since(ti)) / float64(c.Timeout)
	now := sum(c.ToInts(cpu_filter))
	c.Close()

	stat := float64(now - before)
	return Percent(((stat / (td * 100)) * 100) / ncpu)
}
