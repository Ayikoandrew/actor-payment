package receivers

import (
	"fmt"

	"github.com/Ayikoandrew/ap/msg/msg"
	"github.com/Ayikoandrew/ap/storage"
	"github.com/anthdm/hollywood/actor"
)

type PrivateProcessor struct {
	databaseProcessorPID *actor.PID
}

func NewPrivateProcessor() actor.Producer {
	return func() actor.Receiver {
		return &PrivateProcessor{}
	}
}

func (p *PrivateProcessor) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		p.OnStart(ctx)
	case actor.Initialized:
		p.OnInit(ctx)
	case actor.Stopped:
		p.OnStop(ctx)
	case *msg.PrivatePayment:
		ctx.Forward(p.databaseProcessorPID)
	default:
		fmt.Printf("[Private Processor] Unknown message%+v\n", m)
	}
}

func (p *PrivateProcessor) OnStart(ctx *actor.Context) {
	fmt.Println("[Private Processor] Started")
	p.databaseProcessorPID = ctx.SpawnChild(storage.NewDatabaseProcessor(), "database")
}

func (p *PrivateProcessor) OnInit(ctx *actor.Context) {
	fmt.Println("[Private Processor] OnInit")
}

func (p *PrivateProcessor) OnStop(ctx *actor.Context) {
	fmt.Println("[Private Processor] Stopped")
}
