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
var chatID *int

func process(broker *tg.Broker, update tg.APIMessage) {
	// Process messages from marked chat only
	if update.Chat.ChatID != *chatID {
		return
	}
	getNick(update.User)
	updateStats(update)
}

func main() {
	brokerAddr := flag.String("broker", "localhost:7314", "Broker address:port")
	boltdbFile := flag.String("boltdb", "stats.db", "BoltDB database file")
	chatID = flag.Int("chatid", -14625256, "Telegram Chat ID to count stats for")
	flag.Parse()

	var err error
	db, err = bolt.Open(*boltdbFile, 0600, nil)
	assert(err)
	defer db.Close()

	loadUsers()
	loadStats()

	err = tg.CreateBrokerClient(*brokerAddr, process)
	assert(err)
}
