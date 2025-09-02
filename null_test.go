package date

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"
)

func TestNullDate_MarshalJSON(t *testing.T) {
	cases := []struct {
		date NullDate
		want []byte
	}{
		{
			date: NullDate{
				Valid: false,
			},
			want: nullBytes,
		},
		{
			date: NullDate{
				Valid: true,
				Date:  Date{2022, 3, 1},
			},
			want: []byte("\"2022-03-01\""),
		},
	}

	for _, c := range cases {
		t.Run(string(c.want), func(t *testing.T) {
			got, err := c.date.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON(): %v", err)
			}

			if !bytes.Equal(got, c.want) {
				t.Errorf("MarshalJSON() = %v, <nil>; want %v, <nil>", got, c.want)
			}
		})
	}
}

func TestNullDate_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input []byte
		want  NullDate
	}{
		{
			input: nullBytes,
			want:  NullDate{},
		},
		{
			input: []byte("\"2023-03-01\""),
			want: NullDate{
				Valid: true,
				Date:  Date{2023, 3, 1},
			},
		},
	}

	for _, c := range cases {
		t.Run(string(c.input), func(t *testing.T) {
			var date NullDate
			err := date.UnmarshalJSON(c.input)
			if err != nil {
				t.Fatalf("UnmarshalJSON(): %v", err)
			}

			if date.Valid != c.want.Valid || !date.Date.Equal(c.want.Date) {
				t.Errorf("NullDate = %v; want %v", date, c.want)
			}
		})
	}
}

func TestNullDate_UnmarshalJSON_Errors(t *testing.T) {
	cases := [][]byte{
		{},
		[]byte("\"\""),
	}

	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			var date NullDate
			err := date.UnmarshalJSON(c)
			if err == nil {
				t.Error("UnmarshalJSON(v) = <nil>; want error")
			}
		})
	}
}

func TestNullDate_Scan(t *testing.T) {
	cases := []struct {
		input interface{}
		want  NullDate
	}{
		{
			input: nil,
			want: NullDate{
				Valid: false,
			},
		},
		{
			input: time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC),
			want: NullDate{
				Valid: true,
				Date:  Date{2023, 9, 28},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.input), func(t *testing.T) {
			var date NullDate
			err := date.Scan(c.input)
			if err != nil {
				t.Fatalf("Scan(%v): %v", c.input, err)
			}

			if date.Valid != c.want.Valid || !date.Date.Equal(c.want.Date) {
				t.Errorf("NullDate = %v; want %v", date.Valid, c.want.Valid)
			}
		})
	}
}

func TestNullDate_Value(t *testing.T) {
	cases := []struct {
		date NullDate
		want driver.Value
	}{
		{
			date: NullDate{Valid: false},
			want: nil,
		},
		{
			date: NullDate{
				Valid: true,
				Date:  Date{2023, 10, 18},
			},
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
