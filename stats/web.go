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

const USAGE_THRESHOLD = 3

func webWords(rw http.ResponseWriter, req *http.Request) {
	// Filter words under a certain usage
	filtered := make(map[string]UserCount)
	for word, usage := range words {
		total := 0
		for _, count := range usage {
			total += count
		}

		if total < USAGE_THRESHOLD {
			continue
		}
		filtered[word] = usage
	}

	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(filtered)
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
