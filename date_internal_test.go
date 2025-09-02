package date

import (
	"fmt"
	"strconv"
	"testing"
)

func TestDate_add(t *testing.T) {
	cases := []struct {
		date                Date
		years, months, days int
		want                string
	}{
		{Date{2023, 3, 15}, 1, 0, 0, "2024-03-15"},
		{Date{2023, 3, 15}, 0, 1, 0, "2023-04-15"},
		{Date{2023, 3, 15}, 0, 0, 1, "2023-03-16"},
		{Date{2004, 2, 29}, 1, 0, 0, "2005-03-01"},
		{Date{2004, 2, 29}, 0, 1, 0, "2004-03-29"},
		{Date{2004, 2, 29}, 0, 0, 1, "2004-03-01"},
		{Date{2004, 2, 29}, 4, 0, 0, "2008-02-29"},
		{Date{2023, 8, 24}, 20, 5, 3, "2044-01-27"},
		{Date{2023, 8, 24}, 2, 30, 15, "2028-03-10"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s+%d %d %d", c.date.String(), c.years, c.months, c.days), func(t *testing.T) {
			got := c.date.add(c.years, c.months, c.days)
			if got.String() != c.want {
				t.Errorf("%v.add(%d, %d, %d) = %v; want %v", c.date, c.years, c.months, c.days, got, c.want)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	leapYears := []int{4, 1600, 2000, 2004, 2008, 2012, 2016, 2020, 2024, 2028}
	nonLeapYears := []int{1700, 1900, 1999, 2001, 2013, 2018, 2019, 2021, 2022, 2025, 2027, 2100}

	for _, y := range leapYears {
		t.Run(strconv.Itoa(y), func(t *testing.T) {
			if !isLeapYear(y) {
				t.Errorf("%d should be a leap year", y)
			}
		})
	}

	for _, y := range nonLeapYears {
		t.Run(strconv.Itoa(y), func(t *testing.T) {
			if isLeapYear(y) {
				t.Errorf("%d should not be a leap year", y)
			}
		})
	}
}
