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
	BaseURL     string /* Base URL for webhook */
	WebhookURL  string /* Webhook URL */
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

var api Telegram

func main() {
	cfgpath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	file, err := os.Open(*cfgpath)
	assert(err)

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	assert(err)

	// Create Telegram API object
	api = mkAPI(config.Token)

	// Setup webhook handler
	go func() {
		http.HandlerFunc(config.Token, webhook)
		err := http.ListenAndServe(config.BindServer, nil)
		assert(err)
	}()

	// Register webhook @ Telegram
	api.setWebhook(config.BaseURL + config.WebhookURL)

	// Create server for clients
	startClientsServer(config.BindClients)
}
