package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"flag"
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
	destinationOpts := mqtt.NewClientOptions().AddBroker(*destinationBrokerString).SetClientID("mqtt_sync")

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

	sourceClient := mqtt.NewClient(sourceOpts)
	destinationClient := mqtt.NewClient(destinationOpts)

	defer sourceClient.Disconnect(10)
	defer destinationClient.Disconnect(10)

	if token := sourceClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := destinationClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var republishCallback = func(c mqtt.Client, message mqtt.Message) {
		destinationTopic := *destinationTopicPrefix + message.Topic()

		if *debugMode {
			fmt.Printf("%s%s %s => %s%s %s\n",
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

	if token := sourceClient.Subscribe(*sourceTopic, 0, republishCallback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	for true {

	}

}