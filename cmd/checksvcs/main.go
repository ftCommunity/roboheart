package main

import (
	"github.com/ftCommunity-roboheart/roboheart/internal/servicemanager"
	"log"
)

func main() {
	log.Println("Checking dependencies")
	if _, err := servicemanager.NewServiceManager(); err != nil {
		log.Fatal(err)
	}
	log.Println("Dependencies ok!")
}
