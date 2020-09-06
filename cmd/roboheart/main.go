package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ftCommunity-roboheart/roboheart/internal/servicemanager"
)

func main() {
	log.Println("Starting roboheart")
	//create ServiceManager
	sm, err := servicemanager.NewServiceManager()
	if err != nil {
		log.Fatal(err)
	}
	//initialize ServiceManager
	sm.Init()
	log.Println("Start-up completed")
	//setup ctrl-c interrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//wait for ctrl-c
	<-c
	//initiate stop procedure
	log.Println("Stopping roboheart")
	sm.Stop()
	log.Println("Heart rate zero")
	log.Println("Dead")
}
