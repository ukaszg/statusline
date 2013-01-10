package status

import (
	"errors"
	"fmt"
	"strings"
)

type BatteryStatus int
const (
	CHARGING BatteryStatus = iota
	DISCHARGING
	MISSING
	UNKNOWN
)
var (
	full_files = []string{`charge_full`, `energy_full`}
	now_files = []string{`charge_now`, `energy_now`}
	rate_files = []string{`current_now`, `power_now`}
	present_files = []string{`present`}
)
type Battery struct {
	stat map[string]*StatFile
	variant string
	dir_path string
}
func NewBattery(dir_path string, by_design bool) *Battery {
	b := new(Battery)
	if by_design {
		b.variant = "_design" 
	} else {
		b.variant = ""
	}
	b.dir_path = dir_path
	b.stat = make(map[string]*StatFile)

	return b
}
func (b *Battery) getval(cache_name string, filenames []string,  postfix string) (int, error) {
	var (
		stat *StatFile
		ok bool
	)
	if stat, ok = b.stat[cache_name]; ok {
		ok = stat.Open()
	}
	if !ok {
		for _, file_name := range filenames {
			stat = NewStatFile(b.dir_path + file_name + postfix)
			if ok = stat.Open(); ok  {
				break
			}
		}
		if ok {
			b.stat[cache_name] = stat
		} else {
			delete(b.stat, cache_name)
			return -1, errors.New(
				fmt.Sprintf("No files were present:[%s] %s [%s]", b.dir_path, filenames, postfix))
		}
	}
	defer stat.Close()
	return stat.ToInts(func(s []string) []string { return s[0:1] })[0], nil
}
func (b *Battery) Present() (bool, error) {
	if present, err := b.getval("present", present_files, ""); err != nil {
		return false, err
	} else {
		return (present == 1), nil
	}
	panic("unreachable")
}
func (b *Battery) Full() (int, error) {
	return b.getval("full", full_files, b.variant)
}
func (b *Battery) Now() (int, error) {
	return b.getval("now", now_files, "")
}
func (b *Battery) Rate() (int, error) {
	return b.getval("rate", rate_files, "")
}
// Current battery state
func (b *Battery) Get() Percent {
	is_present, present_err := b.Present()
	now, err_now := b.Now()
	full, err_full := b.Full()
	if !is_present || present_err != nil || err_now != nil || err_full != nil {
		return -1.0
	}
	return Percent((float64(now) / float64(full)) * 100.0)
}
func (b *Battery) Status() BatteryStatus {
	var (
		stat *StatFile
		ok bool
	)
	if p, e := b.Present(); e != nil || !p {
		return MISSING
	}
	if stat, ok = b.stat["status"]; !ok {
		stat = NewStatFile(b.dir_path + "status")
	}
	if stat.Open() {
		defer stat.Close()
		b.stat["status"] = stat
		switch strings.TrimSpace(stat.Line()) {
		case "Discharging":
			return DISCHARGING
		case "Charging": fallthrough
		case "Unknown":
			return CHARGING
		}
	}
	return UNKNOWN
}
