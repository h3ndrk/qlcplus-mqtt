package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/NIPE-SYSTEMS/qlcplus-mqtt/api"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type adaptor struct {
	mqttClient      mqtt.Client
	websocketClient *websocket.Conn

	websocketRecvDone chan struct{}
}

func NewAdaptor(mqttOpts *mqtt.ClientOptions, websocketURL string) (*adaptor, error) {
	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	websocketClient, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		return nil, err
	}

	websocketRecvDone := make(chan struct{})

	a := &adaptor{
		mqttClient:        mqttClient,
		websocketClient:   websocketClient,
		websocketRecvDone: websocketRecvDone,
	}

	go func() {
		defer close(websocketRecvDone)
		for {
			t, msg, err := websocketClient.ReadMessage()
			if err != nil {
				log.Print(errors.Wrap(err, "Failed to read websocket message"))
				return
			}
			if t != websocket.TextMessage {
				log.Printf("Unexpected websocket message type: %d", t)
				continue
			}
			log.Printf("recv: %s", msg)

			msgParts := strings.SplitN(string(msg), "|", 1)
			if len(msgParts) < 2 {
				log.Print("Ignoring websocket message: missing separator")
				continue
			}
			if msgParts[0] != "QLC+API" {
				log.Print("Ignoring websocket message: unknown type")
				continue
			}

			go a.handleAPI(msgParts[1])
		}
	}()

	// GOROUTINES:
	// publisher
	// receiver
	// ticker

	return a, nil
}

func (a *adaptor) Stop() {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := a.websocketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to send close message"))
	}
	select {
	case <-a.websocketRecvDone:
	case <-time.After(time.Second):
	}

	a.websocketClient.Close()
	a.mqttClient.Disconnect(250)
}

func (a *adaptor) handleAPI(msg string) {
	msgParts := strings.Split(string(msg), "|")
	if len(msgParts) < 1 {
		log.Printf("Unexpected length of message parts: %d", len(msgParts))
		return
	}

	switch msgParts[0] {
	case "getFunctionsList":
	case "getFunctionType":
	}
}

func (a *adaptor) PublishFunctions() error {
	err := a.websocketClient.WriteMessage(websocket.TextMessage, []byte("QLC+API|getFunctionsList"))
	if err != nil {
		return errors.Wrap(err, "Failed to get list of functions")
	}

	t, msg, err := a.websocketClient.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "Failed to read list of functions")
	}
	if t != websocket.TextMessage {
		return errors.Errorf("Unexpected websocket message type: %d", t)
	}
	msgParts := strings.Split(string(msg), "|")
	if len(msgParts) < 2 {
		return errors.Errorf("Unexpected length of websocket message parts: %d", len(msgParts))
	}

	return nil
}

func tickingPublisher(client mqtt.Client) func() {
	ticker := time.NewTicker(10 * time.Second)
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
	a, err := api.NewAPI("ws://localhost:9999/qlcplusWS")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	f, err := a.GetFunctionsList()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Functions: %+v", f)

	loaded, err := a.IsProjectLoaded()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded: %t", loaded)

	number, err := a.GetFunctionsNumber()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Number: %d", number)

	/*c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	// opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	cx, _, err := websocket.DefaultDialer.Dial("ws://localhost:1234", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer cx.Close()

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

	log.Printf("Got signal %v, exiting.", <-c)*/
}
