package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Ayikoandrew/ap/receivers"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

func main() {
	AccountID := flag.Int64("broker", 0, "broker id")
	flag.Parse()

	r := remote.New("127.0.0.1:40000", remote.NewConfig())
	config := actor.NewEngineConfig().WithRemote(r)

	engine, err := actor.NewEngine(config)
	if err != nil {
		fmt.Printf("Error starting a new engine %+v", err)
	}

	brokerId := engine.Spawn(
		receivers.NewBroker(*AccountID),
		"broker",
		actor.WithID(fmt.Sprintf("broker_%v\n", *AccountID)),
		actor.WithContext(context.TODO()),
		actor.WithInboxSize(1024*2),
	)
	log.Printf("Started broker with PID: %s", brokerId)

	select {}

}
