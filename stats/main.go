package main

import (
	"flag"

	"github.com/hamcha/clessy/tg"
)

func process(broker *tg.Broker, update tg.APIMessage) {
}

func main() {
	brokerAddr := flag.String("broker", "localhost:7314", "Broker address:port")
	flag.Parse()

	err := tg.CreateBrokerClient(*brokerAddr, process)
	if err != nil {
		panic(err)
	}
}
