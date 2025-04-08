package receivers

import (
	"fmt"

	"github.com/Ayikoandrew/ap/msg/msg"
	"github.com/Ayikoandrew/ap/storage"
	"github.com/anthdm/hollywood/actor"
)

type PublicProcessor struct {
	databaseProcessorPID *actor.PID
}

func NewPublicProcessor() actor.Producer {
	return func() actor.Receiver {
		return &PublicProcessor{}
	}
}

func (p *PublicProcessor) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		p.OnStart(ctx)
	case actor.Stopped:
		p.OnStop(ctx)
	case actor.Initialized:
		p.OnInit(ctx)
	case *msg.PublicPayment:
		ctx.Forward(p.databaseProcessorPID)
	default:
		fmt.Printf("[Public Processor] Unknown message%+v\n", m)

	}
}

func (p *PublicProcessor) OnStart(ctx *actor.Context) {
	fmt.Println("[Public Processor] Started")
	p.databaseProcessorPID = ctx.SpawnChild(storage.NewDatabaseProcessor(), "database")
}

func (p *PublicProcessor) OnStop(ctx *actor.Context) {
	fmt.Println("[Public Processor] Stopped")
}

func (p *PublicProcessor) OnInit(ctx *actor.Context) {
	fmt.Println("[Public Processor] Initialized")
}
