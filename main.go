package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

	if token := client.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Printf("Got signal %v, exiting.", <-c)
}
