package main

import (
	"io"
	"strings"
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
			want: "Frijoles (1) caduca mañana",
		},

		test{
			food: Food{Name: "Frijoles", Qty: 1, When: time.Now().AddDate(0, 0, 4)},
			want: "Frijoles (1) caduca en 4 días",
		},

		test{
			food: Food{Name: "Molida", Qty: 1, When: time.Now().AddDate(0, 0, -2)},
			want: "Molida (1) caducó hace 2 días",
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
		if got := daysBetween(tt.a, tt.b); got != tt.want {
			t.Errorf("got=[%d], want=[%d]", got, tt.want)
		}
	}
}

func Test_shouldIgnoreLine(t *testing.T) {
	type test struct {
		line string
		want bool
	}

	tests := []test{
		{line: "hola", want: false},
		{line: "# comment", want: true},
		{line: "", want: true},
	}

	for _, tt := range tests {
		if got := shouldIgnoreLine(tt.line); got != tt.want {
			t.Errorf("got=[%t], want=[%t]", got, tt.want)
		}
	}
}

func Test_extractFoodFromText(t *testing.T) {
	type test struct {
		text       string
		shouldFail bool
		want       Food
	}

	tests := []test{
		test{
			text:       "Frijoles, 1, 2020-05-23",
			want:       Food{Name: "Frijoles", Qty: 1, When: time.Date(2020, 5, 23, 0, 0, 0, 0, time.UTC)},
			shouldFail: false,
		},
		test{
			text:       "Frijoles, 1",
			want:       Food{},
			shouldFail: true,
		},
		test{
			text:       "Frijoles, 1, 2020-",
			want:       Food{},
			shouldFail: true,
		},
		test{
			text:       "Frijoles, x, 2020-05-23",
			want:       Food{},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		got, err := extractFoodFromText(tt.text)
		if tt.shouldFail && (err == nil) {
			t.Errorf("It shouldn't have failed, it failed with: %s error", err.Error())
		}
		if got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_readFile(t *testing.T) {
	type test struct {
		file       io.Reader
		want       []Food
		shouldFail bool
	}

	tests := []test{
		test{
			file: strings.NewReader(`
# Refri:
Molida, 1, 2020-11-23
Milanesas, 2, 2020-11-24
`),
			want: []Food{
				Food{Name: "Molida", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "Milanesas", Qty: 2, When: time.Date(2020, 11, 24, 0, 0, 0, 0, time.UTC)},
			},
			shouldFail: false,
		},

		test{
			file: strings.NewReader(`
# Refri:
Molida, 1, 2020-`),
			want:       []Food{},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		got, err := readFile(tt.file)
		if tt.shouldFail && (err == nil) {
			t.Errorf("It shouldn't have failed, it failed with: %s error", err.Error())
		}
		if !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_equal(t *testing.T) {
	type test struct {
		a, b     []Food
		areEqual bool
	}

	tests := []test{
		test{
			a: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			b: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			areEqual: true,
		},
		test{
			a: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			b: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			areEqual: false,
		},
		test{
			a: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			b: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 2, When: time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)},
			},
			areEqual: false,
		},
	}

	for _, tt := range tests {
		if got := equal(tt.a, tt.b); got != tt.areEqual {
			t.Errorf("got=[%t], want=[%t]", got, tt.areEqual)
		}
	}
}

func Test_filterFoodsByMinDays(t *testing.T) {
	type test struct {
		foods   []Food
		today   time.Time
		minDays int
		want    []Food
	}

	tests := []test{
		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
				Food{Name: "d", Qty: 1, When: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)},
				Food{Name: "e", Qty: 1, When: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)},
			},
			today:   time.Date(2020, 11, 4, 0, 0, 0, 0, time.UTC),
			minDays: 8,
			want: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		if got := filterFoodsByMinDays(&tt.foods, tt.minDays, tt.today); !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_foodByExpiricyDays_Len(t *testing.T) {
	type test struct {
		foods foodByExpiricyDays
		want  int
	}

	tests := []test{
		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		if len := tt.foods.Len(); len != tt.want {
			t.Errorf("got=[%d], want=[%d]", len, tt.want)
		}
	}
}

func Test_foodByExpiricyDays_Swap(t *testing.T) {
	type test struct {
		foods, want foodByExpiricyDays
	}

	tests := []test{
		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
			},
			want: []Food{
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		tt.foods.Swap(0, 1)
		if !equal(tt.foods, tt.want) {
			t.Errorf("got=[%s], want=[%s]", tt.foods, tt.want)
		}
	}
}

func Test_foodByExpiricyDays_Less(t *testing.T) {
	type test struct {
		foods        foodByExpiricyDays
		shouldBeLess bool
	}

	tests := []test{
		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
			},
			shouldBeLess: true,
		},

		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
			},
			shouldBeLess: false,
		},
	}

	for _, tt := range tests {
		if got := tt.foods.Less(0, 1); got != tt.shouldBeLess {
			t.Errorf("got=[%t], want=[%t]", got, tt.shouldBeLess)
		}
	}
}

func Test_alreadyExpired(t *testing.T) {
	type test struct {
		foods []Food
		date  time.Time
		want  []Food
	}

	tests := []test{
		test{
			foods: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
				Food{Name: "d", Qty: 1, When: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)},
				Food{Name: "e", Qty: 1, When: time.Date(2020, 11, 14, 0, 0, 0, 0, time.UTC)},
				Food{Name: "f", Qty: 1, When: time.Date(2020, 11, 15, 0, 0, 0, 0, time.UTC)},
			},
			date: time.Date(2020, 11, 14, 0, 0, 0, 0, time.UTC),
			want: []Food{
				Food{Name: "a", Qty: 1, When: time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC)},
				Food{Name: "b", Qty: 2, When: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)},
				Food{Name: "c", Qty: 1, When: time.Date(2020, 11, 12, 0, 0, 0, 0, time.UTC)},
				Food{Name: "d", Qty: 1, When: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		if got := alreadyExpired(&tt.foods, tt.date); !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}
