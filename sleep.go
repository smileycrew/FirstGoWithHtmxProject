package main

import (
	"fmt"
	"time"
)

func sleep(duration time.Duration) {
	fmt.Printf("Sleeping for %s seconds.", duration)

	var seconds = duration * time.Second

	time.Sleep(seconds)

	fmt.Println("Done sleeping!")
}