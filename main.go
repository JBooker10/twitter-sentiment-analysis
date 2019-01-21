package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cdipaolo/sentiment"
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

// GetRecentTrends returns the most recent trends dependent upon WOIED (Where On Earth ID) specifics
func GetRecentTrends(w http.ResponseWriter, r *http.Request) {
	auth := TwitterAuth()
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	vars := mux.Vars(r)
	searchkey, ok := vars["query"]

	fmt.Println(searchkey)

	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search}`)
		return
	}

	places := &twitter.TrendsPlaceParams{
		WOEID: 1,
	}

	trend, _, err := client.Trends.Place(1, places)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(trend)
}

type UserTimeline struct {
	ID              string
	Name            string
	ScreenName      string
	ProfileImageURL string
	CreatedAt       string
	Location        string
	Text            string
	RetweetCount    int
	FavoriteCount   int
	SentimentRating uint8
}

// GetUserTimeline returns the most recent tweets of a given user with Sentiment Analysis
func GetUserTimeline(w http.ResponseWriter, r *http.Request) {
	auth := TwitterAuth()
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
	}

	vars := mux.Vars(r)
	searchkey, ok := vars["query"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search}`)
		return
	}

	userTimeline := &twitter.UserTimelineParams{
		Count:      10,
		ScreenName: searchkey,
	}

	timeline, _, err := client.Timelines.UserTimeline(userTimeline)
	if err != nil {
		log.Fatal(err)
	}

	tweets := []*UserTimeline{}

	// fmt.Println(timeline)
	for _, tweet := range timeline { // 0
		userTweet := new(UserTimeline)
		userTweet.ID = tweet.IDStr
		userTweet.Name = tweet.User.Name
		userTweet.ScreenName = tweet.User.ScreenName
		userTweet.CreatedAt = tweet.CreatedAt
		userTweet.Location = tweet.User.Location
		userTweet.Text = tweet.Text
		userTweet.RetweetCount = tweet.RetweetCount
		userTweet.FavoriteCount = tweet.FavoriteCount
		userTweet.ProfileImageURL = tweet.User.ProfileImageURL
		userTweet.SentimentRating = model.SentimentAnalysis(tweet.Text, sentiment.English).Score
		tweets = append(tweets, userTweet)
	}

	// fmt.Printf("%+v\n", &tweets)

	json.NewEncoder(w).Encode(tweets)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := mux.NewRouter()
	r.HandleFunc("/search-tweets/{query}", SearchTweets).Methods("GET")
	r.HandleFunc("/search-users/{query}", SearchUsers).Methods("GET")
	r.HandleFunc("/show-trends/{query}", GetRecentTrends).Methods("GET")
	r.HandleFunc("/show-timeline/{query}", GetUserTimeline).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
