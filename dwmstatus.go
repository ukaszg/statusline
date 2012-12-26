package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	STAT    = "/proc/stat"
	MEM     = "/proc/meminfo"
	LAVG    = "/proc/loadavg"
	MEM_VAL = 1
)

/* marks functions to be deferred */
type FileCloser func()

func chkerr(err error) {
	if (err != nil) {
		panic(err.Error())
	}
}

func openFile(file_path string) (*bufio.Reader, FileCloser) {
	var file *os.File
	var err error

	file, err = os.Open(file_path)
	chkerr(err)
	return bufio.NewReader(file), func() { file.Close() }
}

func getInt(r *bufio.Reader, token int) int {
	var i int
	var err error
	var line string
	line, err = r.ReadString('\n')
	chkerr(err)
	i, err = strconv.Atoi(strings.Fields(line)[token])
	chkerr(err)
	return i
}

func memory() string {
	r, closer := openFile(MEM)
	defer closer()

	memTotal := getInt(r, MEM_VAL)
	memFree := memTotal
	for i := 1; i <= 3; i++ {
		memFree -= getInt(r, MEM_VAL)
	}
	return fmt.Sprintf("%d/%dMB", memFree/1024, memTotal/1024)
}

func main() {
	fmt.Println(memory())
}
