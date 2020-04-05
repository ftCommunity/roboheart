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
	log.Println("Waiting 5s until shutdown")
	time.Sleep(5 * time.Second)
	log.Println("Stopping roboheart")
	sm.Stop()
	log.Println("End")
}
