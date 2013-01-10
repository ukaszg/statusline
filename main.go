package main


import (
	"./status"
	"fmt"
	"time"
)

var (
	cpu     = status.NewCpu(1 * time.Second)
	mem     = status.NewMem()
	lavg    = status.NewLoadavg()
	bat     = status.NewBattery("/sys/class/power_supply/BAT1/", false)
)

func Time() string {
	return time.Now().Local().Format("Mon 2 Jan 15:04")
}

func main() {
	c  := cpu.Get()
	cs := status.Bar(16, c)
	m  := mem.Get()
	ms := status.Bar(16, m)
	l  := lavg.Get()
	t  := Time()
	b  := bat.Get()

	fmt.Printf("mem: %s %.1f%% | cpu: %s %.1f%%  %s | bat: %.0f%% | %s\n", ms, m, cs, c, l, b, t)
}
