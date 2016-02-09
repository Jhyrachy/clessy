package tg

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Broker struct {
	Socket net.Conn
}

func ConnectToBroker(brokerAddr string) (*Broker, error) {
	sock, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, err
	}

	broker := new(Broker)
	broker.Socket = sock
	return broker, nil
}

func (b *Broker) Close() {
	b.Socket.Close()
}

func (b *Broker) SendTextMessage(chat *APIChat, text string) {
	cmd := ClientCommand{
		Type: CmdSendTextMessage,
		TextMessageData: &ClientTextMessageData{
			Text: text,
		},
	}
	// Encode command and send to broker
	err := json.NewEncoder(b.Socket).Encode(&cmd)
	if err != nil {
		log.Printf("[SendTextMessage] JSON Encode error: %s\n", err.Error())
	}
	fmt.Fprintf(b.Socket, "\n")
}
