package receivers

import (
	"fmt"

	"github.com/Ayikoandrew/ap/msg/msg"
	"github.com/anthdm/hollywood/actor"
)

type Processor struct {
	publicProcessorPID  *actor.PID
	privateProcessorPID *actor.PID
}

func NewProcessor() actor.Producer {
	return func() actor.Receiver {
		return &Processor{}
	}
}

func (p *Processor) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		p.OnStart(ctx)
	case *msg.PublicPayment:
		ctx.Forward(p.publicProcessorPID)

	case *msg.PrivatePayment:
		ctx.Forward(p.privateProcessorPID)
	default:
		fmt.Printf("Unknown stream %+v\n", m)
	}
}

func (p *Processor) OnStart(ctx *actor.Context) {
	fmt.Println("[Processor Started]")
	p.publicProcessorPID = ctx.SpawnChild(NewPublicProcessor(), "public")
	p.privateProcessorPID = ctx.Engine().Spawn(NewPrivateProcessor(), "private")
}

func (p *Processor) OnStop(ctx *actor.Context) {
	fmt.Println("[Processor Stopped]")
}

func (p *Processor) OnInit(ctx *actor.Context) {
	fmt.Println("[Processor Initialized]")
}
