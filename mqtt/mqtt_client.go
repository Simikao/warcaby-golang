package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"warcaby/game"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func InitMQTT() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("warcaby_server")
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT connection error:", token.Error())
	}
	log.Println("MQTT client connected.")
}

func PublishGameUpdate(g *game.Game) {
	topic := fmt.Sprintf("warcaby/game/%d/move", g.ID)
	payload, err := json.Marshal(g)
	if err != nil {
		log.Println("Błąd serializacji stanu gry:", err)
		return
	}
	token := mqttClient.Publish(topic, 0, false, payload)
	token.Wait()
}
