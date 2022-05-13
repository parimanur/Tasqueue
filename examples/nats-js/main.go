package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/kalbhor/tasqueue"
	nats_broker "github.com/kalbhor/tasqueue/brokers/nats-js"
	"github.com/kalbhor/tasqueue/examples/tasks"
	nats_result "github.com/kalbhor/tasqueue/results/nats-js"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	brkr, err := nats_broker.New(nats_broker.Options{
		URL:         "localhost:4222",
		EnabledAuth: false,
		Streams: map[string][]string{
			"default": []string{tasqueue.DefaultQueue},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := nats_result.New(nats_result.Options{
		URL:         "localhost:4222",
		EnabledAuth: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	srv := tasqueue.NewServer(brkr, res)

	srv.RegisterHandler("add", tasks.SumProcessor)

	var chain []*tasqueue.Task

	for i := 0; i < 3; i++ {
		b, _ := json.Marshal(tasks.SumPayload{Arg1: i, Arg2: 4})
		task, err := tasqueue.NewTask("add", b)
		if err != nil {
			log.Fatal(err)
		}
		chain = append(chain, task)
	}

	t, _ := tasqueue.NewChain(chain...)
	srv.Enqueue(ctx, t)

	srv.Start(ctx, tasqueue.Concurrency(5))

	// Create a task payload.
	fmt.Println("exit..")
}