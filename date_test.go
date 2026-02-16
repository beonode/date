package date_test

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"

	. "github.com/beonode/date"
)

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestFromISO8601(t *testing.T) {
	cases := []struct {
		iso8601 string
		want    Date
	}{
		{"2023-08-15", Date{2023, 8, 15}},
		{"2001-12-24", Date{2001, 12, 24}},
		{"2004-02-29", Date{2004, 2, 29}},
	}

	for _, c := range cases {
		t.Run(c.iso8601, func(t *testing.T) {
			res, err := FromISO8601(c.iso8601)
			if err != nil {
				t.Fatalf("FromISO8601(%s): %v", c.iso8601, err)
			}

			if !res.Equal(c.want) {
				t.Errorf("FromISO8601(%s) = %v; want %v", c.iso8601, res, c.want)
			}
		})
	}
}

func TestFromISO8601_Errors(t *testing.T) {
	cases := []string{"2023-08-32", "2003-02-29"}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			d, err := FromISO8601(c)
			if err == nil {
				t.Errorf("FromISO8601(%s) = %v, %v; want Date{}, error", c, d, err)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestFromTime(t *testing.T) {
	ti := time.Date(2024, 5, 9, 12, 0, 0, 0, time.UTC)
	got := FromTime(ti)
	want := Date{2024, 5, 9}
	if !got.Equal(want) {
		t.Errorf("FromTime(%v) = %v; want %v", ti, got, want)
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestNew(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		day   int
		want  Date
	}{
		{2023, 8, 24, Date{2023, 8, 24}},
		{0, 1, 1, Date{0, 1, 1}},
		{2004, 2, 29, Date{2004, 2, 29}},
		{2005, 6, 16, Date{2005, 6, 16}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%d %d %d", c.year, c.month, c.day), func(t *testing.T) {
			got, err := New(c.year, c.month, c.day)
			if err != nil {
				t.Fatalf("New(%v, %v, %v): %v", c.year, c.month, c.day, err)
			}

			if !got.Equal(c.want) {
				t.Errorf("New(%v, %v, %v) = %v; want %v", c.year, c.month, c.year, got, c.want)
			}
		})
	}
}

func TestNew_Errors(t *testing.T) {
	cases := []struct {
		name  string
		year  int
		month time.Month
		day   int
	}{
		{"invalid month (too low)", 2023, 0, 15},
		{"invalid month (too high)", 2023, 13, 15},
		{"invalid day (too low)", 2023, 5, 0},
		{"invalid day (too high for non leap year)", 2003, 2, 29},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := New(c.year, c.month, c.day)
			if err == nil {
				t.Errorf("New(%d, %d, %d) = %v, %v; want Date{}, error", c.year, c.month, c.day, d, err)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_MarshalJSON(t *testing.T) {
	cases := []struct {
		date Date
		want []byte
	}{
		{Date{2022, 3, 1}, []byte("\"2022-03-01\"")},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			json, err := c.date.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON(): %v", err)
			}

			if !bytes.Equal(json, c.want) {
				t.Errorf("MarshalJSON() = %v, <nil>; want %v, <nil>", json, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		data []byte
		want Date
	}{
		{[]byte("\"2022-03-01\""), Date{2022, 3, 1}},
	}

	for _, c := range cases {
		t.Run(string(c.data), func(t *testing.T) {
			var date Date
			err := date.UnmarshalJSON(c.data)
			if err != nil {
				t.Fatalf("UnmarshalJSON(%v): %v", c.data, err)
			}

			if !date.Equal(c.want) {
				t.Errorf("date = %v; want %v", date, c.want)
			}
		})
	}
}

func TestDate_UnmarshalJSON_Errors(t *testing.T) {
	cases := [][]byte{
		{},
	}

	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			var date Date
			err := date.UnmarshalJSON(c)
			if err == nil {
				t.Error("UnmarshalJSON(v) = <nil>; want error")
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_String(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{2023, 8, 15}, "2023-08-15"},
		{Date{2001, 12, 24}, "2001-12-24"},
		{Date{2004, 2, 29}, "2004-02-29"},
	}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			got := c.date.String()
			if got != c.want {
				t.Errorf("%#v.String() = %s; want %s", c.date, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_ShortString(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{1, 2, 29}, "010229"},
		{Date{2023, 8, 15}, "230815"},
		{Date{2001, 12, 24}, "011224"},
		{Date{2004, 2, 29}, "040229"},
		{Date{4892, 2, 29}, "920229"},
	}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			got := c.date.ShortString()
			if got != c.want {
				t.Errorf("%#v.String() = %s; want %s", c.date, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_IsBefore(t *testing.T) {
	cases := []struct {
		a, b Date
		want bool
	}{
		{Date{2023, 8, 25}, Date{2023, 8, 29}, true},
		{Date{2023, 8, 29}, Date{2023, 8, 25}, false},
		{Date{2023, 8, 25}, Date{2023, 7, 25}, false},
		{Date{2023, 8, 25}, Date{2023, 9, 25}, true},
		{Date{2023, 8, 25}, Date{2022, 8, 25}, false},
		{Date{2023, 8, 25}, Date{2024, 8, 25}, true},
		{Date{2023, 8, 29}, Date{2023, 8, 29}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s is before %s", c.a.String(), c.b.String()), func(t *testing.T) {
			got := c.a.IsBefore(c.b)
			if got != c.want {
				t.Errorf("%v.IsBefore(%v) = %v; want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_IsAfter(t *testing.T) {
	cases := []struct {
		a, b Date
		want bool
	}{
		{Date{2023, 8, 25}, Date{2023, 8, 29}, false},
		{Date{2023, 8, 29}, Date{2023, 8, 25}, true},
		{Date{2023, 8, 25}, Date{2023, 7, 25}, true},
		{Date{2023, 8, 25}, Date{2023, 9, 25}, false},
		{Date{2023, 8, 25}, Date{2022, 8, 25}, true},
		{Date{2023, 8, 25}, Date{2024, 8, 25}, false},
		{Date{2023, 8, 29}, Date{2023, 8, 29}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s is after %s", c.a.String(), c.b.String()), func(t *testing.T) {
			got := c.a.IsAfter(c.b)
			if got != c.want {
				t.Errorf("%v.IsAfter(%v) = %v; want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_Time(t *testing.T) {
	d := Date{2024, 5, 9}
	got := d.Time(time.UTC)
	want := time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC)
	if got != want {
		t.Errorf("%v.Time(%s) = %v; want %v", d, time.UTC, got, want)
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_AddDays(t *testing.T) {
	cases := []struct {
		date Date
		days int
		want Date
	}{
		{
			Date{2024, time.May, 28},
			2,
			Date{2024, time.May, 30},
		},
		{
			Date{2024, time.February, 28},
			1,
			Date{2024, time.February, 29},
		},
		{
			Date{2023, time.February, 28},
			1,
			Date{2023, time.March, 1},
		},
		{
			Date{2023, time.July, 15},
			31,
			Date{2023, time.August, 15},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s+%dd", c.date.String(), c.days), func(t *testing.T) {
			got := c.date.AddDays(c.days)
			if !got.Equal(c.want) {
				t.Errorf("%v.AddDays(%v) = %v; want %v", c.date, c.days, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_AddMonths(t *testing.T) {
	cases := []struct {
		date   Date
		months int
		want   string
	}{
		{Date{2023, 3, 15}, 1, "2023-04-15"},
		{Date{2023, 12, 15}, 12, "2024-12-15"},
		{Date{2004, 2, 29}, 12, "2005-03-01"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s+%d", c.date.String(), c.months), func(t *testing.T) {
			got := c.date.AddMonths(c.months)
			if got.String() != c.want {
				t.Errorf("%v.AddMonths(%v) = %v; want %v", c.date, c.months, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_AddYears(t *testing.T) {
	cases := []struct {
		date  Date
		years int
		want  string
	}{
		{Date{2023, 3, 15}, 1, "2024-03-15"},
		{Date{2004, 2, 29}, 3, "2007-03-01"},
		{Date{2004, 2, 29}, 4, "2008-02-29"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s+%d", c.date.String(), c.years), func(t *testing.T) {
			got := c.date.AddYears(c.years)
			if got.String() != c.want {
				t.Errorf("%v.AddYears(%d) = %v; want %v", c.date, c.years, got, c.want)
			}
		})
	}
}

func TestFirstOfWeek(t *testing.T) {

}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_FirstOfWeek(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{2023, 8, 21}, "2023-08-21"},
		{Date{2023, 8, 22}, "2023-08-21"},
		{Date{2023, 8, 23}, "2023-08-21"},
		{Date{2023, 8, 24}, "2023-08-21"},
		{Date{2023, 8, 25}, "2023-08-21"},
		{Date{2023, 8, 26}, "2023-08-21"},
		{Date{2023, 8, 27}, "2023-08-21"},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.FirstOfWeek()
			if got.String() != c.want {
				t.Errorf("%v.FirstOfWeek() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

func TestStartOfMonth(t *testing.T) {
	sthlm, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		time time.Time
		want time.Time
	}{
		{
			time.Date(2024, 10, 31, 10, 54, 0, 73, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 2, 10, 23, 59, 59, 1e9-1, sthlm),
			time.Date(2024, 2, 1, 0, 0, 0, 0, sthlm),
		},
	}

	for _, c := range cases {
		t.Run(c.time.String(), func(t *testing.T) {
			got := StartOfMonth(c.time)
			if !got.Equal(c.want) {
				t.Errorf("StartOfMonth(%v) = %v; want %v", c.time, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_FirstOfMonth(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{2023, 8, 28}, "2023-08-01"},
		{Date{2023, 9, 30}, "2023-09-01"},
		{Date{2004, 2, 29}, "2004-02-01"},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.FirstOfMonth()
			if got.String() != c.want {
				t.Errorf("%v.FirstOfMonth() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	sthlm, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		time time.Time
		want time.Time
	}{
		{
			time.Date(2024, 10, 31, 10, 54, 0, 73, time.UTC),
			time.Date(2024, 10, 31, 23, 59, 59, 1e9-1, time.UTC),
		},
		{
			time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 29, 23, 59, 59, 1e9-1, time.UTC),
		},
		{
			time.Date(2024, 2, 10, 23, 59, 59, 1e9-1, sthlm),
			time.Date(2024, 2, 29, 23, 59, 59, 1e9-1, sthlm),
		},
	}

	for _, c := range cases {
		t.Run(c.time.String(), func(t *testing.T) {
			got := EndOfMonth(c.time)
			if !got.Equal(c.want) {
				t.Errorf("EndOfMonth(%v) = %v; want %v", c.time, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_LastOfMonth(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{2004, 2, 10}, "2004-02-29"},
		{Date{2003, 2, 10}, "2003-02-28"},
		{Date{2023, 8, 28}, "2023-08-31"},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.LastOfMonth()
			if got.String() != c.want {
				t.Errorf("%v.LastOfMonth() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_Equal(t *testing.T) {
	cases := []struct {
		a, b Date
		want bool
	}{
		{Date{2023, 8, 24}, Date{2023, 8, 24}, true},
		{Date{2023, 8, 24}, Date{2023, 7, 24}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s=%s", c.a.String(), c.b.String()), func(t *testing.T) {
			got := c.a.Equal(c.b)
			if got != c.want {
				t.Errorf("%v.Equal(%v) = %v; want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_Compare(t *testing.T) {
	cases := []struct {
		a, b Date
		want int
	}{
		{Date{2023, 8, 24}, Date{2023, 8, 24}, 0},
		{Date{2023, 8, 24}, Date{2023, 8, 25}, -1},
		{Date{2023, 8, 24}, Date{2023, 8, 23}, 1},
		{Date{2023, 8, 24}, Date{2023, 9, 24}, -1},
		{Date{2023, 8, 24}, Date{2023, 7, 24}, 1},
		{Date{2023, 8, 24}, Date{2024, 8, 24}, -1},
		{Date{2023, 8, 24}, Date{2022, 8, 24}, 1},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s cmp %s", c.a.String(), c.b.String()), func(t *testing.T) {
			got := c.a.Compare(c.b)
			if got != c.want {
				t.Errorf("%v.Compare(%v) = %v; want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_LastOfWeek(t *testing.T) {
	cases := []struct {
		date Date
		want string
	}{
		{Date{2023, 8, 21}, "2023-08-27"},
		{Date{2023, 8, 22}, "2023-08-27"},
		{Date{2023, 8, 23}, "2023-08-27"},
		{Date{2023, 8, 24}, "2023-08-27"},
		{Date{2023, 8, 25}, "2023-08-27"},
		{Date{2023, 8, 26}, "2023-08-27"},
		{Date{2023, 8, 27}, "2023-08-27"},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.LastOfWeek()
			if got.String() != c.want {
				t.Errorf("%v.LastOfWeek() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_Scan(t *testing.T) {
	cases := []struct {
		input interface{}
		want  Date
	}{
		{
			input: time.Date(2023, time.September, 28, 0, 0, 0, 0, time.UTC),
			want:  Date{2023, time.September, 28},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.input), func(t *testing.T) {
			var date Date
			err := date.Scan(c.input)
			if err != nil {
				t.Fatalf("Scan(%v): %v", c.input, err)
			}

			if !date.Equal(c.want) {
				t.Errorf("Date = %v; want %v", date, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_Value(t *testing.T) {
	cases := []struct {
		date Date
		want driver.Value
	}{
		{
			date: Date{2023, 10, 18},
			want: "2023-10-18",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.date), func(t *testing.T) {
			got, err := c.date.Value()
			if err != nil {
				t.Fatalf("Value(): %v", err)
			}

			if got != c.want {
				t.Errorf("Value() = %v, <nil>; want %v, <nil>", got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestEndOfDay(t *testing.T) {
	date := Date{2023, 8, 24}
	endOfDay := date.EndOfDay(time.UTC)
	want := time.Date(2023, time.Month(8), 24, 23, 59, 59, 1e9-1, time.UTC)

	if !endOfDay.Equal(want) {
		t.Errorf("EndOfDay(%v, UTC) = %v; want %v", date, endOfDay, want)
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestStartOfDay(t *testing.T) {
	day := time.Date(2023, time.August, 24, 13, 37, 0, 0, time.UTC)
	startOfDay := StartOfDay(day)
	want := time.Date(2023, time.August, 24, 0, 0, 0, 0, time.UTC)

	if !startOfDay.Equal(want) {
		t.Errorf("StartOfDay(%v) = %v; want %v", day, startOfDay, want)
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_StartOfDay(t *testing.T) {
	date := Date{2023, 8, 24}
	got := date.StartOfDay(time.UTC)
	want := time.Date(2023, time.Month(8), 24, 0, 0, 0, 0, time.UTC)

	if !got.Equal(want) {
		t.Errorf("%v.StartOfDay(UTC) = %v; want %v", date, got, want)
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDiffInDays(t *testing.T) {
	cases := []struct {
		d1   Date
		d2   Date
		want int
	}{
		{
			d1:   Date{2023, 9, 1},
			d2:   Date{2023, 9, 2},
			want: 1,
		},
		{
			d1:   Date{2023, 9, 2},
			d2:   Date{2023, 9, 1},
			want: 1,
		},
		{
			d1:   Date{2023, 9, 1},
			d2:   Date{2023, 9, 1},
			want: 0,
		},
		{
			d1:   Date{2023, 8, 1},
			d2:   Date{2023, 9, 1},
			want: 31,
		},
		{
			d1:   Date{2004, 2, 28},
			d2:   Date{2004, 3, 1},
			want: 2,
		},
		{
			d1:   Date{2004, 2, 28},
			d2:   Date{2005, 2, 28},
			want: 366,
		},
		{
			d1:   Date{2005, 2, 28},
			d2:   Date{2006, 2, 28},
			want: 365,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s - %s", c.d1.String(), c.d2.String()), func(t *testing.T) {
			got := DiffInDays(c.d1, c.d2)
			if got != c.want {
				t.Errorf("DiffInDays(%v, %v) = %v; want %v", c.d1, c.d2, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_StartOfMonth(t *testing.T) {
	cases := []struct {
		date Date
		want Date
	}{
		{
			Date{2023, time.February, 15},
			Date{2023, time.February, 1},
		},
		{
			Date{2024, time.May, 31},
			Date{2024, time.May, 1},
		},
		{
			Date{2024, time.January, 1},
			Date{2024, time.January, 1},
		},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.StartOfMonth()
			if !got.Equal(c.want) {
				t.Errorf("%v.StartOfMonth() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

//goland:noinspection GoStructInitializationWithoutFieldNames
func TestDate_EndOfMonth(t *testing.T) {
	cases := []struct {
		date Date
		want Date
	}{
		{
			Date{2023, time.February, 15},
			Date{2023, time.February, 28},
		},
		{
			Date{2024, time.February, 15},
			Date{2024, time.February, 29},
		},
		{
			Date{2024, time.May, 31},
			Date{2024, time.May, 31},
		},
		{
			Date{2024, time.January, 1},
			Date{2024, time.January, 31},
		},
	}

	for _, c := range cases {
		t.Run(c.date.String(), func(t *testing.T) {
			got := c.date.EndOfMonth()
			if !got.Equal(c.want) {
				t.Errorf("%v.StartOfMonth() = %v; want %v", c.date, got, c.want)
			}
		})
	}
}

var benchAdditions = []int{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
}

func BenchmarkDate_AddDays(b *testing.B) {
	start := Date{
		Year:  2023,
		Month: 8,
		Day:   24,
	}

	for _, a := range benchAdditions {
		b.Run(strconv.Itoa(a), func(b *testing.B) {
			var date Date
			for i := 0; i < b.N; i++ {
				date = start.AddDays(a)
			}
			runtime.KeepAlive(date)
		})
	}
}

func BenchmarkDate_AddMonths(b *testing.B) {
	start := Date{
		Year:  2023,
		Month: 8,
		Day:   24,
	}

	for _, a := range benchAdditions {
		b.Run(strconv.Itoa(a), func(b *testing.B) {
			var date Date
			for i := 0; i < b.N; i++ {
				date = start.AddMonths(a)
			}
			runtime.KeepAlive(date)
		})
	}
}

func BenchmarkDate_AddYears(b *testing.B) {
	start := Date{
		Year:  2023,
		Month: 8,
		Day:   24,
	}

	for _, a := range benchAdditions {
		b.Run(strconv.Itoa(a), func(b *testing.B) {
			var date Date
			for i := 0; i < b.N; i++ {
				date = start.AddYears(a)
			}
			runtime.KeepAlive(date)
		})
	}
}
