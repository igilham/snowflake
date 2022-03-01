package main

import (
	"fmt"
	"log"
	"time"

	"github.com/igilham/snowflake"
)

const (
	// EPOCH_YEAR is the year of the epoch as defined for our use case.
	EPOCH_YEAR = 2020

	// EPOCH_MONTH is the month of the epoch as defined for our use case.
	EPOCH_MONTH = 1

	// EPOCH_DAY is the month of the epoch as defined for our use case.
	EPOCH_DAY = 1
)

func main() {
	epoch := time.Date(EPOCH_YEAR, EPOCH_MONTH, EPOCH_DAY, 0, 0, 0, 0, time.UTC).Unix()
	fmt.Printf("epoch: %v\n\n", epoch)
	var id uint16 = 1001
	gen := snowflake.NewWorker(epoch, id)
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())

	fmt.Printf("\n...\n\n")
	dur, err := time.ParseDuration("1s")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(dur)

	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
	fmt.Printf("%v\n", gen.Next())
}
