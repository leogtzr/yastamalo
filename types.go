package main

import "time"

// Food ...
type Food struct {
	Name string
	Qty  int
	When time.Time
}

// ExpirationFormatter ...
type ExpirationFormatter interface {
	Format() string
}

type foodByExpiricyDays []Food
