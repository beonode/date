# date
Date representation in Go.

## Installation
```sh
go get -u github.com/beonode/date
```

## Example usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/beonode/date"
)

func main() {
	// Create a date (error handling omitted).
	d, _ := date.New(2025, 12, 10)
	d, _ = date.FromISO8601("2025-12-10")
	d = date.FromTime(time.Now())
	d = date.Today(time.UTC)

	// Get relative dates.
	startOfMonth := d.StartOfMonth()
	lastOfWeek := d.LastOfWeek()

	// Or go back to time.Time.
	startOfDay := d.StartOfDay(time.UTC)
	endOfDay := d.EndOfDay(time.UTC)

	fmt.Println(d, startOfMonth, lastOfWeek, startOfDay, endOfDay)
}
```

Both `Date` and `NullDate` implement the `sql.Scanner` and `driver.Valuer` interfaces.