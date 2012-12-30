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

func Memory() string {
	r, closer := open_file(MEM)
	defer closer()

	memTotal := get_mem_stat(r)
	memFree := memTotal
	for i := 1; i <= 3; i++ {
		memFree -= get_mem_stat(r)
	}
	return fmt.Sprintf("%d/%dMB", memFree/1024, memTotal/1024)
}

var (
	BarStart string = "["
	BarEnd   string = "]"
	BarEmpty string = " "
	BarFull  string = "|"
)

func Bar(inside_length int, value float64) string {
	ret := make([]string, inside_length + len(BarStart) + len(BarEnd))
	till := int( math.Ceil(( float64(inside_length) / 100.0 ) * value ) )

	ret = append(ret, BarStart)
	for i := 0; i <= till; i++ {
		ret = append(ret, BarFull)
	}
	for i := till; i < inside_length; i++ {
		ret = append(ret, BarEmpty)
	}
	ret = append(ret, BarEnd)

	return strings.Join(ret, "") 
}

func main() {
	//fmt.Println(Memory())
	//fmt.Printf("%.1f%%", Cpu())
	cpu := Cpu()
	fmt.Printf("mem: %s | cpu: %s %.1f%%", Memory(), Bar(15, cpu), cpu )
}
