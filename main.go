package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

var interrupted bool

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting...")

	nc.Subscribe("insights.entries.info.updated", func(m *nats.Msg) {
		fmt.Println(fmt.Sprintf("Recieved via push %s", string(m.Data)))

		if m.Reply != "" {
			m.Respond([]byte(""))
		}
	})

	Pull(nc, "$JS.API.CONSUMER.MSG.NEXT.ENTRIES.UPDATED_PULL", 2*time.Second, func(m *nats.Msg) {
		fmt.Println(fmt.Sprintf("Recieved via pull %s", string(m.Data)))
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for range c {
		fmt.Println("Received SIGINT, cleaning up...")
		interrupted = true
		nc.Close()
		return
	}
}

// Pull like Subscribe, but for pulling
func Pull(nc *nats.Conn, topic string, timeout time.Duration, handle nats.MsgHandler) {
	go func() {
		for {
			if interrupted == true {
				return
			}
			msg, err := nc.Request(topic, []byte(""), timeout)
			if err == nats.ErrTimeout {
				continue
			}
			if err != nil {
				fmt.Print(fmt.Sprintf("Error: %s", err))
				continue
			}

			handle(msg)

			msg.Respond(nil)
		}
	}()
}
