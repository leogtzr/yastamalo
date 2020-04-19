package main

import (
	"fmt"
	"time"
)

func (f Food) String() string {
	now := time.Now()
	daysBetweenNowAndFoodExpiry := daysBetween(now, f.When)
	if daysBetweenNowAndFoodExpiry == 0 {
		return fmt.Sprintf("%s (%d) caduca hoy", f.Name, f.Qty)
	}
	return fmt.Sprintf("%s (%d) caduca en %d días", f.Name, f.Qty, daysBetweenNowAndFoodExpiry)
}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}
