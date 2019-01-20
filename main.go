package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// SearchTweets - Searches All Tweets
func SearchTweets(w http.ResponseWriter, r *http.Request) {

	auth := TwitterAuth()
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	vars := mux.Vars(r)
	searchkey, ok := vars["query"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search}`)
		return
	}

	searchTweets := &twitter.SearchTweetParams{
		Query:           searchkey,
		ResultType:      "recent",
		Count:           2,
		IncludeEntities: twitter.Bool(true),
	}

	search, _, err := client.Search.Tweets(searchTweets)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(search)
}

// SearchUsers - Searches All Users
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	auth := TwitterAuth()
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	vars := mux.Vars(r)
	searchkey, ok := vars["query"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search}`)
		return
	}

	searchUsers := &twitter.UserSearchParams{
		Count:           4,
		IncludeEntities: twitter.Bool(true),
	}

	search, _, err := client.Users.Search(searchkey, searchUsers)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(search)
}

func ShowTrends(w http.ResponseWriter, r *http.Request) {
	auth := TwitterAuth
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	getTrends := 
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := mux.NewRouter()
	r.HandleFunc("/search-tweets/{query}", SearchTweets).Methods("GET")
	r.HandleFunc("/search-users/{query}", SearchUsers).Methods("GET")
	r.HandleFunc("/show-trends/{query}", ShowTrends).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
