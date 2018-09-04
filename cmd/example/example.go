package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kodek/sense-api"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	client, err := sense_api.NewClient(os.Getenv("SENSE_USER"), os.Getenv("SENSE_PASSWORD"))
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	recv, err := client.Realtime(done)
	if err != nil {
		panic(err)
	}

	go func() {
		throttle := time.Tick(1 * time.Second)
		for val := range recv {
			<-throttle
			fmt.Printf("Consumption [%.2f] Production [%.2f]\n", val.Payload.W, val.Payload.SolarW)
		}
	}()

	time.Sleep(5 * time.Second)
	done <- struct{}{}
	time.Sleep(2 * time.Second)
}
