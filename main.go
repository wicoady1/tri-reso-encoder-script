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

	configPath := "config.ini"
	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config, err: %v", err)
	}

	var stop func()
	defer func() {
		if stop != nil {
			stop()
		}
		log.Println("Stopping...")
	}()

	op := operations.NewOperations(conf)
	go func() {
		for true {
			op.EncodeVideoInBackground()
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

}
