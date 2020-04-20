package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	dbFile := flag.String("db", "", "db file")
	flag.Parse()
	fmt.Println(*dbFile)

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
	filteredFoods := filterFoodsByMinDays(&foods, minDaysToAlertUserForExpiry, today)

	sort.Sort(foodByExpiricyDays(filteredFoods))
	for _, food := range filteredFoods {
		fmt.Println(food)
	}
}
