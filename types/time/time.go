// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

// time is built on top of the go standard time package
// provides additional functions for time

package time

import (
	"gosimple/types"
	"math"
	"time"
)

// Is evaluates whether 'a' is time.Time
func Is(a any) bool {
	return types.IsTime(a)
}

// From converts param 'a' of a basic type to time.Time
// Returns error if param 'a' type is not:
//   string, int, float, uint or time
func From(a any) (time.Time, error) {
	return types.ToTime(a)
}

// FromString converts a numeric string to time.Time
// Similar to time.Parse(format, s)
// Returns error if param 's' type is not string
// or can't be converted to time
func FromString(s any) (time.Time, error) {
	return types.StringToTime(s)
}

// FromInt converts any int type representing unix time to time.Time
// Equivilant to time.Unix(i, 0)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func FromInt(i any) (time.Time, error) {
	return types.IntToTime(i)
}

// FromFloat converts any float type representing unix time to time.Time
// Sumilar to time.Unix(int64(i), n)
// Returns error if param 'f' type is not float32, float64
func FromFloat(f any) (time.Time, error) {
	return types.FloatToTime(f)
}

// FromUint converts any uint type representing unix time to time.Time
// Equivilant to time.Unix(int64(i), 0)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func FromUint(u any) (time.Time, error) {
	return types.UintToTime(u)
}

// FromTime converts a time.Time to asserted time.Time
// Returns error if param 't' type is not time.Time
func FromTime(t any) (time.Time, error) {
	return types.TimeToTime(t)
}

// String converts a time to string
// Equivilant to fmt.Sprint(t)
// Returns error if param 't' type is not time.Time
func String(t any) (string, error) {
	return types.TimeToString(t)
}

// Int converts a time.Time to a unix int
// Equivilant to int(t.Unix())
// Returns error if param 't' type is not time.Time
func Int(t any) (int, error) {
	return types.TimeToInt(t)
}

// Float converts a time.Time to a unix float64
// Equivilant to float64(t.Unix())
// Returns error if param 't' type is not time.Time
func Float(t any) (float64, error) {
	return types.TimeToFloat(t)
}

// Uint converts a time.Time to a unix uint
// Equivilant to uint(t.Unix())
// Returns error if param 't' type is not time.Time
func Uint(t any) (uint, error) {
	return types.TimeToUint(t)
}

// DaysInMonth returns the number of calendar days
// for year 'y' and month 'm' provided
func DaysInMonth(y int, m int) int {
	return time.Date(y, time.Month(m), 0, 0, 0, 0, 0, time.Now().Location()).Day()
}

// MonthStart returns the first date of the month for time 't'
func MonthStart(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// MonthEnd returns the last nanosecond of the month for time 't'
func MonthEnd(t time.Time) time.Time {
	return MonthStart(t).AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
}

// QuarterStart returns the first date of the quarter
// for time 't' with year ending in month 'ye'
func QuarterStart(t time.Time, ye time.Month) time.Time {
	ye = ye % 3
	ye = (3-((t.Month()-ye)%3))%3 - 2
	return MonthStart(t.AddDate(0, int(ye), 0))
}

// QuarterEnd returns the last nanosecond of the quarter
// for time 't' with year ending in month 'ye'
func QuarterEnd(t time.Time, ye time.Month) time.Time {
	ye = ye % 3
	ye = (3 - ((t.Month() - ye) % 3)) % 3
	return MonthEnd(t.AddDate(0, int(ye), 0))
}

// YearStart returns the first date of the year
// for time 't' with year ending in month 'ye'
func YearStart(t time.Time, ye time.Month) time.Time {
	y := 0
	if t.Month() < ye+1 {
		y = 1
	}
	return time.Date(t.Year()-y, ye+1, 1, 0, 0, 0, 0, t.Location())
}

// YearEnd returns the last nanosecond of the year
// for time 't' with year ending in month 'ye'
func YearEnd(t time.Time, ye time.Month) time.Time {
	return YearStart(t, ye).AddDate(0, 12, 0).Add(-1 * time.Nanosecond)
}

// HOLIDAYS
// methods and storage for standard and custom holidays

type Holiday struct {
	time.Time
	// Name is the common name of the holiday
	Name string
	// Date returns the date of the holiday for year 'y'
	Date func(y int) time.Time
}

type Holidays struct {
	List []Holiday
}

func GetUsHolidays() Holidays {
	h := Holidays{}
	h.List = append(h.List, Holiday{Name: "New Years Day", Date: NewYears})
	h.List = append(h.List, Holiday{Name: "Martin Luther King Day", Date: MlkDay})
	h.List = append(h.List, Holiday{Name: "Inauguration Day", Date: InagurationDay})
	h.List = append(h.List, Holiday{Name: "Presidents Day", Date: PresidentsDay})
	h.List = append(h.List, Holiday{Name: "Memorial Day", Date: MemorialDay})
	h.List = append(h.List, Holiday{Name: "Juneteenth", Date: NationalIndependenceDay})
	h.List = append(h.List, Holiday{Name: "Independence Day", Date: IndependenceDay})
	h.List = append(h.List, Holiday{Name: "Labor Day", Date: LaborDay})
	h.List = append(h.List, Holiday{Name: "Columbus Day", Date: ColumbusDay})
	h.List = append(h.List, Holiday{Name: "Veterans Day", Date: VeteransDay})
	h.List = append(h.List, Holiday{Name: "Thanksgiving", Date: Thanksgiving})
	h.List = append(h.List, Holiday{Name: "Christmas", Date: Christmas})
	return h
}

func (h *Holidays) IsHoliday(t time.Time) bool {
	y, m, d := t.Date()
	for _, i := range h.List {
		_, hm, hd := i.Date(y).Date()
		if m == hm && d == hd {
			return true
		}
	}
	return false
}

// Instance returns the date of the 'i' instance of weekday 'w'
// in month 'm' of year 'y'; if i < 0 returns the last instance, and
// panics if 'i' is 0 or exceeds the number of instances
func Instance(i int, d time.Weekday, m time.Month, y int) time.Time {
	s := time.Date(y, m, 1, 0, 0, 0, 0, time.Now().Location())
	e := MonthEnd(s).Add(-24*time.Hour + time.Nanosecond)
	f := s.Weekday()
	l := e.Weekday()
	o := 0
	if i < 0 {
		if d > l {
			o = 7
		}
		return e.AddDate(0, 0, int(d-l)-o)
	}
	if d >= f {
		o = 7
	}
	r := s.AddDate(0, 0, i*7+int(d-f)-o)
	if r.After(e) || r.Before(s) {
		panic("gosimple.types.time.Instance: invalid instance of date; must not exceed instances in month and not be 0")
	}
	return r
}

// NewYears returns the observed date for new years day of for year 'y'
func NewYears(y int) time.Time {
	return HolidayObserved(time.Date(y, time.January, 1, 0, 0, 0, 0, time.Now().Location()))
}

// MlkDay returns the date of Martin Luther King Jr Day for year 'y'
func MlkDay(y int) time.Time {
	return Instance(3, time.Monday, time.January, y)
}

// InagurationDay returns the date of the presidential inaguration for year 'y'
func InagurationDay(y int) time.Time {
	//if (y-1)%4 > 0 {
	//	return time.Time{}, fmt.Errorf("gosimple.types.time.InaugurationDay: no inauguration in %d", y)
	//}
	y -= (y - 1) % 4
	return HolidayObserved(time.Date(y, time.January, 20, 0, 0, 0, 0, time.Now().Location()))
}

// PresidentsDay returns the date of President's Day (or Washington's Birthday) for year 'y'
func PresidentsDay(y int) time.Time {
	return Instance(3, time.Monday, time.February, y)
}

// GoodFriday returns the date of good friday for year 'y'
func GoodFriday(y int) time.Time {
	return Easter(y).AddDate(0, 0, -2)
}

// Easter returns the date of easter for year 'y'
func Easter(y int) time.Time {
	var yr, c, n, k, i, j, l, m, d float64
	yr = float64(y)
	c = math.Floor(yr / 100)
	n = yr - 19*math.Floor(yr/19)
	k = math.Floor((c - 17) / 25)
	i = c - math.Floor(c/4) - math.Floor((c-k)/3) + 19*n + 15
	i = i - 30*math.Floor(i/30)
	i = i - math.Floor(i/28)*(1-math.Floor(i/28)*math.Floor(29/(i+1))*math.Floor((21-n)/11))
	j = yr + math.Floor(yr/4) + i + 2 - c + math.Floor(c/4)
	j = j - 7*math.Floor(j/7)
	l = i - j
	m = 3 + math.Floor((l+40)/44)
	d = l + 28 - 31*math.Floor(m/4)
	return time.Date(y, time.Month(m-1), int(d), 0, 0, 0, 0, time.Now().Location())
}

// MemorialDay returns the date of Memorial Day for year 'y'
func MemorialDay(y int) time.Time {
	return Instance(-1, time.Monday, time.May, y)
}

// NationalIndependenceDay returns the observed date for
// Junteenth National Independence Day for year 'y'
func NationalIndependenceDay(y int) time.Time {
	return HolidayObserved(time.Date(y, time.June, 19, 0, 0, 0, 0, time.Now().Location()))
}

// IndependenceDay returns the observed date for US Independence Day for year 'y'
func IndependenceDay(y int) time.Time {
	return HolidayObserved(time.Date(y, time.July, 4, 0, 0, 0, 0, time.Now().Location()))
}

// LaborDay returns the date of Labor Day for year 'y'
func LaborDay(y int) time.Time {
	return Instance(1, time.Monday, time.September, y)
}

// ColumbusDay returns the date of Columbus Day for year 'y'
func ColumbusDay(y int) time.Time {
	return Instance(2, time.Monday, time.October, y)
}

// VeteransDay returns the observed date for Veterans Day for year 'y'
func VeteransDay(y int) time.Time {
	return HolidayObserved(time.Date(y, time.November, 11, 0, 0, 0, 0, time.Now().Location()))
}

// Thanksgiving returns the date of Thanksgiving Day for year 'y'
func Thanksgiving(y int) time.Time {
	return Instance(4, time.Thursday, time.November, y)
}

// Christmas returns the observed date for Christmas Day for year 'y'
func Christmas(y int) time.Time {
	return HolidayObserved(time.Date(y, time.December, 25, 0, 0, 0, 0, time.Now().Location()))
}

// HolidayObserved returns the date holiday 'h' is observed,
// Friday if on Saturday and Monday if on Sunday
func HolidayObserved(h time.Time) time.Time {
	if h.Weekday() == time.Saturday {
		h = h.AddDate(0, 0, -1)
	} else if h.Weekday() == time.Sunday {
		h = h.AddDate(0, 0, 1)
	}
	return h
}
