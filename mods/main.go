package main

import (
	"flag"

	"../tg"
)

func dispatch(broker *tg.Broker, update tg.APIMessage) {
	metafora(broker, update)
}

func main() {
	brokerAddr := flag.String("broker", "localhost:7314", "Broker address:port")
	flag.Parse()

	err := tg.CreateBrokerClient(*brokerAddr, dispatch)
	if err != nil {
		panic(err)
	}
}
