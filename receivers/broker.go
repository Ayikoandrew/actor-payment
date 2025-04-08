package receivers

import (
	"fmt"
	"time"

	"github.com/Ayikoandrew/ap/msg/msg"
	"github.com/anthdm/hollywood/actor"
)

type Broker struct {
	AccountID int64
}

var procID = actor.NewPID("127.0.0.1:30000", "processor/processor_0")

func NewBroker(accountID int64) actor.Producer {
	return func() actor.Receiver {
		return &Broker{
			AccountID: accountID,
		}
	}
}

func (b *Broker) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Initialized:
		b.OnInit(ctx)
	case actor.Started:
		b.OnStart(ctx)
	case actor.Stopped:
		b.OnStop(ctx)
	default:
		fmt.Printf("Unknown message %+v", m)
	}
}

func (b *Broker) OnStart(ctx *actor.Context) {
	fmt.Println("[Broker Started]")

	go b.start(ctx)
}

func (b *Broker) OnStop(ctx *actor.Context) {
	fmt.Println("[Broker Stopped]")
}

func (b *Broker) OnInit(ctx *actor.Context) {
	fmt.Println("[Broker Initialized]")
}

func (b *Broker) start(ctx *actor.Context) {
	fmt.Println("started sending event streams")
	fmt.Printf("processor id %s\n", procID)
	for {
		time.Sleep(2 * time.Second)
		ctx.Send(procID, &msg.PublicPayment{
			AccountID: b.AccountID,
			Amount:    100,
		})

		ctx.Send(procID, &msg.PrivatePayment{
			AccountID: b.AccountID,
			Amount:    1000,
		})
	}
}
