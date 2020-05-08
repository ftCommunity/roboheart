package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Stopping roboheart")
	sm.Stop()
	log.Println("End")
}
