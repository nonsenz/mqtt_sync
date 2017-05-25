package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"flag"
	"os"
	"log"
)

func main() {

	sourceBrokerString := flag.String("s", "tcp://127.0.0.1:1883", "source broker connection string")
	sourceUserString := flag.String("su", "", "source broker username")
	sourcePassString := flag.String("sp", "", "source broker password")
	destinationBrokerString := flag.String("d", "tcp://127.0.0.1:1883", "source broker connection string")
	destinationUserString := flag.String("du", "", "destination broker username")
	destinationPassString := flag.String("dp", "", "destination broker password")
	sourceTopic := flag.String("t", "#", "source topic")
	destinationTopicPrefix := flag.String("p", "", "destination topic prefix (e.g. /foo)")
	debugMode := flag.Bool("debug", false, "turn on debug output")

	flag.Parse()

	sourceOpts := mqtt.NewClientOptions().AddBroker(*sourceBrokerString).SetClientID("mqtt_sync")
	sourceOpts.SetAutoReconnect(true)
	destinationOpts := mqtt.NewClientOptions().AddBroker(*destinationBrokerString).SetClientID("mqtt_sync")
	destinationOpts.SetAutoReconnect(true)

	if *sourceUserString != "" {
		sourceOpts.SetUsername(*sourceUserString)
	}

	if *sourcePassString != "" {
		sourceOpts.SetPassword(*sourcePassString)
	}

	if *destinationUserString != "" {
		destinationOpts.SetUsername(*destinationUserString)
	}

	if *destinationPassString != "" {
		destinationOpts.SetPassword(*destinationPassString)
	}

	destinationClient := mqtt.NewClient(destinationOpts)

	var republishCallback = func(c mqtt.Client, message mqtt.Message) {
		destinationTopic := *destinationTopicPrefix + message.Topic()

		if *debugMode {
			log.Printf("%s%s %s => %s%s %s\n",
				*sourceBrokerString,
				message.Topic(),
				message.Payload(),
				*destinationBrokerString,
				destinationTopic,
				message.Payload())
		}

		token := destinationClient.Publish(destinationTopic, 0, false, message.Payload())
		token.Wait()
	}

	sourceOpts.OnConnect = func(sourceClient mqtt.Client) {
		if token := sourceClient.Subscribe(*sourceTopic, 0, republishCallback); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	sourceClient := mqtt.NewClient(sourceOpts)

	if token := sourceClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("source host: %v\n", token.Error())
		os.Exit(1)
	}

	if token := destinationClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("destination host: %v\n", token.Error())
		os.Exit(1)
	}

	defer sourceClient.Disconnect(10)
	defer destinationClient.Disconnect(10)

	fmt.Println("mqtt_sync connected...")

	for true {

	}

}