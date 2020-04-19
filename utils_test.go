package main

import (
	"testing"
	"time"
)

func TestFood_String(t *testing.T) {
	type test struct {
		food Food
		want string
	}

	tests := []test{
		test{
			food: Food{Name: "Frijoles", Qty: 1, When: time.Now()},
			want: "Frijoles (1) caduca hoy",
		},

		test{
			food: Food{Name: "Frijoles", Qty: 1, When: time.Now().AddDate(0, 0, 1)},
			want: "Frijoles (1) caduca en 1 días",
		},
	}

	for _, tt := range tests {
		if got := tt.food.String(); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_daysBetween(t *testing.T) {
	type test struct {
		a, b time.Time
		want int
	}

	tests := []test{
		test{
			a:    time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		test{
			a:    time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		test{
			a:    time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		test{
			a:    time.Date(2016, 12, 25, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2017, 1, 7, 0, 0, 0, 0, time.UTC),
			want: 13,
		},
	}

	for _, tt := range tests {
		got := daysBetween(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("got=[%d], want=[%d]", got, tt.want)
		}
	}
}
