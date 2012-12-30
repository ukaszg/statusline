package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	STAT    = "/proc/stat"
	MEM     = "/proc/meminfo"
	LAVG    = "/proc/loadavg"
	MEM_VAL = 1
)

var ()

/* marks functions to be deferred */
type FileCloser func()

func chkerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func open_file(file_path string) (*bufio.Reader, FileCloser) {
	var file *os.File
	var err error

	file, err = os.Open(file_path)
	chkerr(err)
	return bufio.NewReader(file), func() { file.Close() }
}

func get_mem_stat(r *bufio.Reader) int {
	var ret int
	line, err := r.ReadString('\n')
	chkerr(err)

	ret, err = strconv.Atoi(strings.Fields(line)[MEM_VAL])
	chkerr(err)

	return ret
}

func get_cpu_stat() (v int) {
	r, closer := open_file(STAT)
	line, err := r.ReadString('\n')
	chkerr(err)
	closer()

	tokens := strings.Fields(line)
	for _, s := range tokens[1:4] {
		i, err := strconv.Atoi(s)
		chkerr(err)
		v += i
	}
	return v
}

var (
	timeout = 1 * time.Second
	ncpu    = float64(runtime.NumCPU())
)

func Cpu() float64 {
	ti := time.Now()
	before := get_cpu_stat()
	time.Sleep(timeout)

	td := float64(time.Since(ti)) / float64(timeout)
	now := get_cpu_stat()

	stat := float64(now - before)
	return ((stat / (td * 100)) * 100) / ncpu
}

func Memory() float64 {
	r, closer := open_file(MEM)
	defer closer()

	memTotal := get_mem_stat(r)
	memFree := memTotal
	for i := 1; i <= 3; i++ {
		memFree -= get_mem_stat(r)
	}
	return float64(memFree) / float64(memTotal) * 100
}

var (
	BarStart string = "["
	BarEnd   string = "]"
	BarEmpty string = " "
	BarFull  string = "|"
)

func Bar(length int, value float64) string {
	ret := make([]string, length)
	inside_length := length - (len(BarStart) + len(BarEnd))
	till := int(math.Ceil( (float64(inside_length) / 100.0) * value ))

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

func Loadavg() string {
	r, closer := open_file(LAVG)
	line, err := r.ReadString('\n')
	chkerr(err)
	closer()
	return strings.Join(strings.Fields(line)[:3], " ")
}

func Time() string {
	return time.Now().Local().Format("Mon 2 Jan 15:04")
}

func main() {
	cpu := Cpu()
	mem := Memory()
	fmt.Printf("mem: %s %.1f%% | cpu: %s %.1f%%  %s | %s\n",
		Bar(16, mem), mem, Bar(16, cpu), cpu, Loadavg(), Time() )
}
