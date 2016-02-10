package main

import (
	"flag"

	"github.com/boltdb/bolt"
	"github.com/hamcha/clessy/tg"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

var db *bolt.DB

func process(broker *tg.Broker, update tg.APIMessage) {

}

func main() {
	brokerAddr := flag.String("broker", "localhost:7314", "Broker address:port")
	boltdbFile := flag.String("boltdb", "stats.db", "BoltDB database file")
	flag.Parse()

	db, err := bolt.Open(*boltdbFile, 0600, nil)
	assert(err)
	defer db.Close()

	loadStats()

	err = tg.CreateBrokerClient(*brokerAddr, process)
	assert(err)
}
