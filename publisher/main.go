package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type UserEvent struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

func main() {
	// Tell the client to talk to the emulator
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "local-dev")
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	topic := client.Topic("user-events")

	event := UserEvent{UserID: "user-123", Action: "login"}
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}

	// Publish is async — it batches messages internally for efficiency.
	// result.Get() blocks until the message is confirmed by the server.
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	msgID, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("result.Get: %v", err)
	}

	fmt.Printf("Published message ID: %s\n", msgID)
}
