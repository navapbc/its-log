package types

// We name things based on dates
func (s *Storage) YYYYMMDD() string {
	return s.ILTime.AsYYYYMMDD()
}

func (s *Storage) SetDateYMD(ymd string) error {
	if ymd == "today" {
		s.ILTime = NewILTimeToday()
		return nil
	}
	d, e := NewILTimeFromYMD(ymd)
	if e == nil {
		s.ILTime = d
	}
	return e
}

func (s *Storage) SetDateILT(ilt *ILTime) {
	s.ILTime = ilt
}

func (s *Storage) SubtractDays(days int) {
	s.ILTime.SubtractDays(days)
}

func (s *Storage) AddDays(days int) {
	s.ILTime.AddDays(days)
}
