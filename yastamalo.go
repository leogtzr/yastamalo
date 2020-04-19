package main

import (
	"flag"
	"fmt"
)

func main() {

	dbFile := flag.String("db", "", "db file")
	flag.Parse()
	fmt.Println(*dbFile)

}
