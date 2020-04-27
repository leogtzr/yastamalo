package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {

	config, err := readConfig("yastamalo.env", os.Getenv("HOME"), map[string]interface{}{
		"minDaysToAlertUserForExpiry": 7,
		"maxDaysToShowInExpiredFoods": 5,
	})

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	dbFile := flag.String("db", "", "db file")
	flag.Parse()

	if len(*dbFile) == 0 {
		panic(fmt.Errorf("db file cannot be empty"))
	}

	file, err := os.Open(*dbFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	foods, err := readFile(file)
	if err != nil {
		panic(err)
	}

	today := time.Now()
	filteredFoods := filterFoodsByMinDays(&foods, config.GetInt("minDaysToAlertUserForExpiry"), today)
	if len(filteredFoods) > 0 {
		sort.Sort(foodByExpiricyDays(filteredFoods))
		fmt.Printf("=== Ya van a caducarrrrrr ===\n")
		for _, food := range filteredFoods {
			fmt.Printf("\t* %s\n", food)
		}
	}

	alreadyExpired := alreadyExpired(&foods, today)
	if len(alreadyExpired) > 0 {
		fmt.Printf("=== Ya caducaron y deberías deshacerte de ellos ===\n")
		for _, food := range foods {
			if food.When.Before(today) {
				fmt.Printf("\t* %s\n", food)
			}
		}
	}
}
