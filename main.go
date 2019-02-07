package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := mux.NewRouter()
	r.HandleFunc("/search-tweets/{query}", SearchTweets).Methods("GET")
	r.HandleFunc("/users/{user}", SearchUsers).Methods("GET")
	r.HandleFunc("/trends/{woeid}", GetRecentTrends).Methods("GET")
	r.HandleFunc("/timeline/{name}", GetUserTimeline).Methods("GET")
	r.HandleFunc("/stream/{query}", StreamUserTweets).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
