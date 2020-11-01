package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ftCommunity-roboheart/roboheart/internal/servicemanager"
)

func main() {
	parser := argparse.NewParser("roboheart", "roboheart")
	configpath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "path to configuration file",
	})
	pluginpaths:=parser.StringList("p","plugin",&argparse.Options{
		Required: false,
		Help:     "path to plugin files",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	log.Println("Starting roboheart")
	var config []byte
	if *configpath != "" {
		var err error
		config, err = ioutil.ReadFile(*configpath)
		if err != nil {
			log.Fatal(err)
		}
	}
	//create ServiceManager
	sm, err := servicemanager.NewServiceManager(config,*pluginpaths)
	if err != nil {
		log.Fatal(err)
	}
	//initialize ServiceManager
	sm.Init()
	log.Println("Start-up completed")
	//setup ctrl-c interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//wait for ctrl-c
	<-c
	//initiate stop procedure
	log.Println("Stopping roboheart")
	sm.Stop()
	log.Println("Heart rate zero")
	log.Println("Dead")
}
