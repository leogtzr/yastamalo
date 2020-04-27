package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func (f Food) String() string {
	now := time.Now()
	daysBetweenNowAndFoodExpiry := daysBetween(now, f.When)
	if daysBetweenNowAndFoodExpiry == 0 {
		return fmt.Sprintf("%s (%d) caduca hoy", f.Name, f.Qty)
	}
	if daysBetweenNowAndFoodExpiry == 1 {
		return fmt.Sprintf("%s (%d) caduca mañana", f.Name, f.Qty)
	}
	if f.When.Before(now) {
		return fmt.Sprintf("%s (%d) caducó hace %d días", f.Name, f.Qty, daysBetweenNowAndFoodExpiry)
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

func shouldIgnoreLine(line string) bool {
	return strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0
}

func readFile(file io.Reader) ([]Food, error) {
	foods := make([]Food, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if shouldIgnoreLine(line) {
			continue
		}
		f, err := extractFoodFromText(line)
		if err != nil {
			return []Food{}, err
		}
		foods = append(foods, f)
	}
	return foods, scanner.Err()
}

func extractFoodFromText(line string) (Food, error) {
	foodFields := strings.Split(line, foodRecordSeparator)
	if len(foodFields) != requiredNumberOfFieldsInFoodRecord {
		return Food{}, fmt.Errorf("error parsing (%s) line, wrong number of fields", line)
	}
	expiryField := strings.TrimSpace(foodFields[2])
	when, err := time.Parse(foodExpiryDateFormatLayout, expiryField)
	if err != nil {
		return Food{}, err
	}
	qty, err := strconv.Atoi(strings.TrimSpace(foodFields[1]))
	if err != nil {
		return Food{}, err
	}
	return Food{Name: foodFields[0], Qty: qty, When: when}, nil
}

func equal(a, b []Food) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func filterFoodsByMinDays(foods *[]Food, minDays int, today time.Time) []Food {
	filteredFoods := make([]Food, 0)
	for _, f := range *foods {

		if f.When.Before(today) {
			continue
		}

		if daysBetween := daysBetween(today, f.When); daysBetween <= minDays {
			filteredFoods = append(filteredFoods, f)
		}
	}
	return filteredFoods
}

// Len ...
func (foods foodByExpiricyDays) Len() int {
	return len(foods)
}

// Swap ...
func (foods foodByExpiricyDays) Swap(i, j int) {
	foods[i], foods[j] = foods[j], foods[i]
}

// Less ...
func (foods foodByExpiricyDays) Less(i, j int) bool {
	food1 := foods[i]
	food2 := foods[j]
	return food1.When.Before(food2.When)
}

func alreadyExpired(foods *[]Food, date time.Time) []Food {
	alreadyExpired := make([]Food, 0)
	for _, food := range *foods {
		if food.When.Before(date) {
			alreadyExpired = append(alreadyExpired, food)
		}
	}
	return alreadyExpired
}

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}
