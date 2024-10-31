package observability

import (
    "log"
    "github.com/DataDog/datadog-go/statsd"
)

var Client *statsd.Client

func Init() {
    var err error
    Client, err = statsd.New("127.0.0.1:8125")
    if err != nil {
        log.Fatalf("Failed to create statsd client: %v", err)
    }
}

func Close() {
    Client.Close()
}
