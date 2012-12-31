package main

import (
	"./status"
	"fmt"
	"time"
)

var (
	timeout = 1 * time.Second
	cpu     = status.NewCpu()
	mem     = status.NewMem()
	lavg    = status.NewLoadavg()
)

func Time() string {
	return time.Now().Local().Format("Mon 2 Jan 15:04")
}

func main() {
	c  := cpu.Get(timeout)
	cs := status.Bar(16, c)
	m  := mem.Get()
	ms := status.Bar(16, m)
	l  := lavg.Get()
	t  := Time()

	fmt.Printf("mem: %s %.1f%% | cpu: %s %.1f%%  %s | %s\n", ms, m, cs, c, l, t)
}
