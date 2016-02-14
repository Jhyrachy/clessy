package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func webStats(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(stats)
	if err != nil {
		log.Println("[webStats] JSON Encoding error: " + err.Error())
	}
}

func webUsers(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(users)
	if err != nil {
		log.Println("[webUsers] JSON Encoding error: " + err.Error())
	}
}

func webWords(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(filteredWords())
	if err != nil {
		log.Println("[webWords] JSON Encoding error: " + err.Error())
	}
}

func startWebServer(bindAddr string) {
	http.HandleFunc("/stats", webStats)
	http.HandleFunc("/users", webUsers)
	http.HandleFunc("/words", webWords)
	http.ListenAndServe(bindAddr, nil)
}
