package main

import (
	"log"
	"time"

	"github.com/ftCommunity/roboheart/internal/servicemanager"
)

func main() {
	log.Println("Starting roboheart")
	sm, err := servicemanager.NewServiceManager()
	if err != nil {
		log.Fatal(err)
	}
	err = sm.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Start-up completed")
}
