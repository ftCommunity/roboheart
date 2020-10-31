package main

import (
	"github.com/ftCommunity-roboheart/roboheart/internal/servicemanager"
	"log"
)

func main() {
	log.Println("Checking services")
	if _, err := servicemanager.NewServiceManager(nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Services ok!")
}
