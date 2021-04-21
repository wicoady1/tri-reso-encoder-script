package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wicoady1/tri-reso-encoder-script/config"
	"github.com/wicoady1/tri-reso-encoder-script/operations"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	conf := config.Config{}

	/*
		conf, err := config.Load(configPath)
		if err != nil {
			log.Fatalf("failed to load config, err: %v", err)
		}
	*/

	/*
		if conf.Pprof.Enabled {
			go func() {
				pprof.Start(conf)
			}()
		}
	*/

	var stop func()
	defer func() {
		if stop != nil {
			stop()
		}
		log.Println("Stopping...")
	}()

	/*
		switch mode {
		case webServiceMode:
			stop = webservice.Start(conf)
		default:
			stop = webservice.Start(conf)
		}
	*/

	go func() {
		for true {
			operations.EncodeVideoInBackground(&conf)
		}
	}()

	log.Println("Working in background. Type \"quit\" to stop")
	//waiting for any key
	for true {
		var test string
		fmt.Scanf("%s", &test)
		log.Println(test)

		if test == "quit" {
			break
		}
	}

	//<-signalChan

}
