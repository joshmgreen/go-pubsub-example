package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"cloud.google.com/go/pubsub"
)

type UserEvent struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

func main() {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client, err := pubsub.NewClient(ctx, "local-dev")
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription("analytics-worker")

	fmt.Println("Waiting for messages... (Ctrl+C to stop)")

	// Receive blocks until ctx is cancelled.
	// It spawns goroutines internally — your handler runs concurrently.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var event UserEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("bad message, nacking: %v", err)
			msg.Nack() // Pub/Sub will redeliver this message
			return
		}

		fmt.Printf("Received: userID=%s action=%s\n", event.UserID, event.Action)

		// Do your work here:
		// - write to ClickHouse
		// - trigger a Temporal workflow
		// - call another service

		msg.Ack() // Tell Pub/Sub: done, don't redeliver
	})

	if err != nil {
		log.Fatalf("Receive: %v", err)
	}
}
