package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a duration for the test service to run like:")
		fmt.Println(os.Args[0], "30s")
		fmt.Println("this will run the test executable for 30 seconds")
		os.Exit(2)
	}
	fmt.Println("run sleep testing process for", os.Args[1])
	iTime, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(iTime)
}
