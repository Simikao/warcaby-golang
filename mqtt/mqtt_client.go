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

func PublishGameWin(gameID int, winnerPiece int) {
	topic := fmt.Sprintf("warcaby/game/%d/win", gameID)
	winner := ""
	switch winnerPiece {
	case 1:
		winner = "Czarne"
	case 2:
		winner = "Białe"
	default:
		winner = "Coś poszło nie tak, nie wiem kto wygrał"
	}

	payloadData := map[string]interface{}{
		"gameID":  gameID,
		"winner":  winner,
		"message": fmt.Sprintf("Grę wygrywa %s", winner),
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		log.Println("Błąd serializacji payloadu o wygranej:", err)
		return
	}

	token := mqttClient.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Println("Błąd publikacji powiadomienia o wygranej:", token.Error())
	} else {
		log.Printf("Opublikowano powiadomienie o wygranej na temacie %s: %s\n", topic, payload)
	}
}
