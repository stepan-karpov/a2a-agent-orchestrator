package main

import (
	"adk"
	"log"
)

func main() {
	entry, err := adk.CreateEntry("eliza")
	if err != nil {
		log.Fatal(err)
	}
	entry.PrintHello()

}
