package status

import (
	"math"
	"strings"
)

var (
	BarStart string = "["
	BarEnd   string = "]"
	BarEmpty string = " "
	BarFull  string = "|"
)

type Percent float64

func Bar(length int, value Percent) string {
	ret := make([]string, length)
	inside_length := length - (len(BarStart) + len(BarEnd))
	till := int(math.Ceil((float64(inside_length) / 100.0) * float64(value)))

	ret = append(ret, BarStart)
	for i := 1; i < length-2; i++ {
		if i <= till {
			ret = append(ret, BarFull)
		} else {
			ret = append(ret, BarEmpty)
		}
	}

	ret = append(ret, BarEnd)

	return strings.Join(ret, "")
}
