package servicemngr

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/servicemngr/core/internal/servicemanager"
	"github.com/servicemngr/core/package/manifest"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(name string, desc string, svcs [][]manifest.ServiceManifest) {
	parser := argparse.NewParser(name, desc)
	configpath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "path to configuration file",
	})
	pluginpaths := parser.StringList("p", "plugin", &argparse.Options{
		Required: false,
		Help:     "path to plugin files",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	log.Println("Starting", name)
	var config []byte
	if *configpath != "" {
		var err error
		config, err = ioutil.ReadFile(*configpath)
		if err != nil {
			log.Fatal(err)
		}
	}
	//create ServiceManager
	sm, err := servicemanager.NewServiceManager(config, *pluginpaths, svcs)
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
	log.Println("Stopping", name)
	sm.Stop()
	log.Println("Heart rate zero")
	log.Println("Dead")
}
