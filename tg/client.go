package tg

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type UpdateHandler func(broker *Broker, message APIMessage)

func CreateBrokerClient(brokerAddr string, updateFn UpdateHandler) error {
	broker, err := ConnectToBroker(brokerAddr)
	if err != nil {
		return err
	}
	defer broker.Close()

	in := bufio.NewReader(broker.Socket)
	for {
		bytes, _, err := in.ReadLine()
		if err != nil {
			break
		}

		var update APIUpdate
		err = json.Unmarshal(bytes, &update)
		if err != nil {
			log.Printf("[tg - CreateBrokerClient] ERROR reading JSON: %s\r\n", err.Error())
			log.Println(string(bytes))
		}

		// Dispatch to UpdateHandler
		updateFn(broker, update.Message)
	}
	return io.EOF
}
