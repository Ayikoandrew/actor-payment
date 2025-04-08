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
	processorID := flag.Int64("processor", 0, "processor id")
	flag.Parse()

	r := remote.New("127.0.0.1:30000", remote.NewConfig())
	config := actor.NewEngineConfig().WithRemote(r)

	engine, err := actor.NewEngine(config)
	if err != nil {
		fmt.Printf("Failed to start processor actor engine %+v", err)
	}

	processId := engine.Spawn(
		receivers.NewProcessor(),
		"processor",
		actor.WithContext(context.Background()),
		actor.WithID(fmt.Sprintf("processor_%d", *processorID)),
		actor.WithInboxSize(1024*2),
	)

	log.Printf("Started broker with PID: %s", processId)

	select {}
}
