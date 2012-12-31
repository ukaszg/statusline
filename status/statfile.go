package status

import (
	"os"
	"strconv"
	"strings"
	"bufio"
)

func chkerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func sum(i []int) (j int) {
	for _, v := range i {
		j += v
	}
	return
}
type FieldFilter func([]string) []string
type StatFile struct {
	r *bufio.Reader
	filename string
}
func(s *StatFile) Open() bool {
	if s.r != nil {
		return false
	}
	if file, err := os.Open(s.filename); err != nil {
		s.r = nil
		return false
	} else {
		s.r = bufio.NewReader(file)
	}
	return true
}
func (s *StatFile) Close() {
	if (s.r == nil) {
		return
	}
	s.r = nil
}
func NewStatFile(file_path string) *StatFile {
	return &StatFile{nil, file_path}
}
func (s *StatFile) Tokens() []string {
	return strings.Fields(s.Line())
}
func (s *StatFile) Line() string {
	line, err := s.r.ReadString('\n')
	chkerr(err)
	return line
}
func (s *StatFile) ToInts(filter FieldFilter) (conv []int) {
	var (
		i int
		e error
	)
	for _, t := range filter(s.Tokens()) {
		i, e = strconv.Atoi(t)
		chkerr(e)
		conv = append(conv, i)
	}
	return
}
