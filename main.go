package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func tickingPublisher(client mqtt.Client) func() {
	ticker := time.NewTicker(time.Second)
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
	loop:
		for {
			select {
			case t := <-ticker.C:
				client.Publish("qlcplus/tick", 0, false, t.String())
			case <-done:
				break loop
			}
		}

		wg.Done()
	}()

	return func() {
		ticker.Stop()
		close(done)
		wg.Wait()
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	// opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer client.Disconnect(250)

	if token := client.Subscribe("go-mqtt/sample", 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Print(string(msg.Payload()))
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	cancel := tickingPublisher(client)
	defer cancel()

	log.Printf("Got signal %v, exiting.", <-c)
}
