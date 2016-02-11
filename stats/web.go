package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func webStats(rw http.ResponseWriter, req *http.Request) {
	err := json.NewEncoder(rw).Encode(stats)
	if err != nil {
		log.Println("[webStats] JSON Encoding error: " + err.Error())
	}
}

func webUsers(rw http.ResponseWriter, req *http.Request) {
	err := json.NewEncoder(rw).Encode(users)
	if err != nil {
		log.Println("[webUsers] JSON Encoding error: " + err.Error())
	}
}

func startWebServer(bindAddr string) {
	http.HandleFunc("/stats", webStats)
	http.HandleFunc("/users", webUsers)
	http.ListenAndServe(bindAddr, nil)
}
