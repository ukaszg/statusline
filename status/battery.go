package status

type Battery struct {
	present   StatFile
	full_1    StatFile
	full_2    StatFile
	now_1     StatFile
	now_2     StatFile
}

func NewBattery(dir_path string, by_design bool) *Battery {
	b := new(Battery)
	var variant string = ""
	if by_design {
		variant = "_design"
	}
	b.full_1 = *NewStatFile(dir_path + "charge_full" + variant)
	b.full_2 = *NewStatFile(dir_path + "energy_full" + variant)
	b.now_1 = *NewStatFile(dir_path + "charge_now")
	b.now_2 = *NewStatFile(dir_path + "energy_now")
	b.present = *NewStatFile(dir_path + "present")
	return b
}

var battery_filter = func(s []string) []string { return s[0:1] }

func (b *Battery) Get() Percent {
	b.present.Open()
	defer b.present.Close()
	if b.present.ToInts(battery_filter)[0] != 1 {
		return NONE
	}
	var now, full int

	if b.full_1.Open() {
		defer b.full_1.Close()
		full = b.full_1.ToInts(battery_filter)[0]
	} else if b.full_2.Open() {
		defer b.full_2.Close()
		full = b.full_2.ToInts(battery_filter)[0]
	} else {
		return NONE
	}

	if b.now_1.Open() {
		defer b.now_1.Close()
		now = b.now_1.ToInts(battery_filter)[0]
	} else if b.now_2.Open() {
		defer b.now_2.Close()
		now = b.now_2.ToInts(battery_filter)[0]
	} else {
		return NONE
	}

	return Percent( (float64(now) / float64(full)) * 100.0 )
}
