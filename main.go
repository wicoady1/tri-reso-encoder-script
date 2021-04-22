package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/wicoady1/tri-reso-encoder-script/config"
	"github.com/wicoady1/tri-reso-encoder-script/operations"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

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

	p := widgets.NewParagraph()
	p.Text = `Triple Resolution Encoder
	----------------------
	1. Drag your mp4 video to folder "input"
	2. Wait until the process is complete
	3. Video will be available in "output" folder`
	p.SetRect(0, 0, 50, 7)
	ui.Render(p)

	typeExit := widgets.NewParagraph()
	typeExit.Text = `Type any key to exit`
	typeExit.SetRect(0, 7, 50, 10)
	ui.Render(typeExit)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}

}
