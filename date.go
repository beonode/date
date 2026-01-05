package date

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func Today(l *time.Location) Date {
	return FromTime(time.Now().In(l))
}

func FromISO8601(date string) (Date, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Date{}, err
	}

	return Date{t.Year(), t.Month(), t.Day()}, nil
}

func FromTime(t time.Time) Date {
	y, m, d := t.Date()
	return Date{y, m, d}
}

func New(year int, month time.Month, day int) (Date, error) {
	if month > 12 || month < 1 {
		return Date{}, fmt.Errorf("month must be between 1-12 (inclusive), got %d", month)
	} else if day < 1 {
		return Date{}, fmt.Errorf("day must be greater than 0, got %d", day)
	}

	lastDayOfMonth := daysInMonth(year, month)
	if day > lastDayOfMonth {
		return Date{}, fmt.Errorf("last day of month %04d-%02d is %d, got %d", year, month, lastDayOfMonth, day)
	}

	return Date{year, month, day}, nil
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) IsBefore(o Date) bool {
	if d.Year < o.Year {
		return true
	}

	if d.Year == o.Year && d.Month < o.Month {
		return true
	}

	if d.Year == o.Year && d.Month == o.Month && d.Day < o.Day {
		return true
	}

	return false
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) IsAfter(o Date) bool {
	if d.Year > o.Year {
		return true
	}

	if d.Year == o.Year && d.Month > o.Month {
		return true
	}

	if d.Year == o.Year && d.Month == o.Month && d.Day > o.Day {
		return true
	}

	return false
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) Time(l *time.Location) time.Time {
	return startOfDay(d.Year, d.Month, d.Day, l)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) AddDays(days int) Date {
	return d.add(0, 0, days)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) AddMonths(months int) Date {
	return d.add(0, months, 0)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) AddYears(years int) Date {
	return d.add(years, 0, 0)
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) FirstOfMonth() Date {
	return Date{d.Year, d.Month, 1}
}

func EndOfMonth(t time.Time) time.Time {
	y, m, d := StartOfMonth(t).Date()
	return time.Date(y, m, d, 23, 59, 59, 1e9-1, t.Location()).AddDate(0, 1, -1)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) LastOfMonth() Date {
	return Date{d.Year, d.Month, daysInMonth(d.Year, d.Month)}
}

func FirstOfWeek(t time.Time) time.Time {
	wd := t.Weekday()
	if wd == time.Monday {
		return t
	} else if wd == time.Sunday {
		return t.AddDate(0, 0, -6)
	}
	return t.AddDate(0, 0, int(-wd)+1)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) FirstOfWeek() Date {
	return FromTime(FirstOfWeek(d.Time(time.UTC)))
}

func LastOfWeek(t time.Time) time.Time {
	wd := t.Weekday()
	if wd == time.Sunday {
		return t
	}
	return t.AddDate(0, 0, 7-int(wd))
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) LastOfWeek() Date {
	return FromTime(LastOfWeek(d.Time(time.UTC)))
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) Equal(o Date) bool {
	return d.Year == o.Year && d.Month == o.Month && d.Day == o.Day
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

//goland:noinspection GoMixedReceiverTypes
func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	date, err := FromISO8601(str)
	if err != nil {
		return err
	}

	d.Year, d.Month, d.Day = date.Year, date.Month, date.Day
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) ShortString() string {
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

//goland:noinspection GoMixedReceiverTypes
func (d *Date) Scan(value any) error {
	if t, ok := value.(time.Time); ok {
		d.Year, d.Month, d.Day = t.Year(), t.Month(), t.Day()
		return nil
	}

	return fmt.Errorf("cannot scan type %T into Date", value)
}

func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return startOfDay(y, m, d, t.Location())
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) StartOfDay(l *time.Location) time.Time {
	return startOfDay(d.Year, d.Month, d.Day, l)
}

func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return endOfDay(y, m, d, t.Location())
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) EndOfDay(l *time.Location) time.Time {
	return endOfDay(d.Year, d.Month, d.Day, l)
}

func DiffInDays(d1 Date, d2 Date) int {
	diff := d1.Time(time.UTC).Sub(d2.Time(time.UTC))
	return int(diff.Abs() / (24 * time.Hour))
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) StartOfMonth() Date {
	return d.add(0, 0, -d.Day+1)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) EndOfMonth() Date {
	return d.add(0, 1, -d.Day)
}

//goland:noinspection GoMixedReceiverTypes
func (d Date) add(years, months, days int) Date {
	return FromTime(d.Time(time.UTC).AddDate(years, months, days))
}

func startOfDay(y int, m time.Month, d int, l *time.Location) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, l)
}

func endOfDay(y int, m time.Month, d int, l *time.Location) time.Time {
	return time.Date(y, m, d, 23, 59, 59, 1e9-1, l)
}

var monthDays = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

func daysInMonth(year int, month time.Month) int {
	if month == 2 && isLeapYear(year) {
		return 29
	}

	return monthDays[int(month)]
}

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}

	if year%100 == 0 {
		return false
	}

	return year%4 == 0
}
