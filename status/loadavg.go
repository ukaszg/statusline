package status

import (
	"strings"
)

type Loadavg struct {
	StatFile
}

func NewLoadavg() *Loadavg {
	return &Loadavg{*NewStatFile("/proc/loadavg")}
}
func (l *Loadavg) Get() string {
	l.Open()
	defer l.Close()
	return strings.Join(l.Tokens()[:3], " ")
}
