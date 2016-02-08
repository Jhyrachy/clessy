package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
)

type Config struct {
	BindServer  string /* Address:Port to bind for Telegram */
	BindClients string /* Address:Port to bind for clients */
	Token       string /* Telegram bot token */
	WebhookURL  string /* Webhook URL */
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfgpath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	file, err := os.Open(*cfgpath)
	assert(err)

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	assert(err)

	// Setup webhook handler
	go func() {
		http.HandlerFunc(config.Token, webhook)
		err := http.ListenAndServe(config.BindServer, nil)
		assert(err)
	}()

	// Create server for clients
	startClientsServer(config.BindClients)
}
