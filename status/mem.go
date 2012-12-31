package status

type Mem struct {
	StatFile
}

func NewMem() *Mem {
	return &Mem{*NewStatFile("/proc/meminfo")}
}

var mem_filter = func(s []string) []string { return s[1:2] }

func (m *Mem) Get() Percent {
	m.Open()
	defer m.Close()

	total := m.ToInts(mem_filter)[0]
	free := total

	for i := 1; i <= 3; i++ {
		free -= m.ToInts(mem_filter)[0]
	}
	return Percent(float64(free) / float64(total) * 100.0)
}
