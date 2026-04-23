# go-pubsub-example

A minimal Go example showing how to publish and consume messages with [Google Cloud Pub/Sub](https://cloud.google.com/pubsub), running entirely against the local emulator — no GCP account required.

## What's in here

| Path | What it does |
|---|---|
| `publisher/main.go` | Publishes a single `UserEvent` JSON message to the `user-events` topic |
| `subscriber/main.go` | Pulls messages from the `analytics-worker` subscription until interrupted (Ctrl+C) |

## Prerequisites

- Go 1.25+
- [Google Cloud Pub/Sub Emulator](https://cloud.google.com/pubsub/docs/emulator) (via `gcloud` or Docker)

## Running locally

**1. Start the emulator**

```bash
gcloud beta emulators pubsub start --project=local-dev --host-port=localhost:8085
```

Or with Docker:

```bash
docker run --rm -p 8085:8085 gcr.io/google.com/cloudsdktool/cloud-sdk \
  gcloud beta emulators pubsub start --project=local-dev --host-port=0.0.0.0:8085
```

**2. Create the topic and subscription**

```bash
export PUBSUB_EMULATOR_HOST=localhost:8085

# Using the Pub/Sub REST API against the emulator
curl -s -X PUT "http://localhost:8085/v1/projects/local-dev/topics/user-events"
curl -s -X PUT "http://localhost:8085/v1/projects/local-dev/subscriptions/analytics-worker" \
  -H "Content-Type: application/json" \
  -d '{"topic":"projects/local-dev/topics/user-events"}'
```

**3. Start the subscriber** (in one terminal)

```bash
go run ./subscriber
```

**4. Publish a message** (in another terminal)

```bash
go run ./publisher
```

You should see the subscriber print:

```
Received: userID=user-123 action=login
```

## How it works

Both programs set `PUBSUB_EMULATOR_HOST=localhost:8085` at startup, which tells the Pub/Sub client library to skip GCP auth and talk directly to the emulator.

The publisher serializes a `UserEvent` struct to JSON, publishes it asynchronously via `topic.Publish`, then calls `result.Get` to block until the server confirms the message.

The subscriber calls `sub.Receive`, which spawns goroutines internally to handle messages concurrently. Each handler unmarshals the JSON, does its work, then calls `msg.Ack()`. On a bad message it calls `msg.Nack()` so Pub/Sub redelivers it.
