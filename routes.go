package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cdipaolo/sentiment"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
)

// SearchTweets - Searches All Tweets
func SearchTweets(w http.ResponseWriter, r *http.Request) {
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

	searchTweets := &twitter.SearchTweetParams{
		Query:           searchkey,
		ResultType:      "recent",
		Lang:            "en",
		Count:           40,
		IncludeEntities: twitter.Bool(false),
	}

	search, _, err := client.Search.Tweets(searchTweets)
	if err != nil {
		log.Fatal(err)
	}

	tweets := []*UserTweets{}

	for _, tweet := range search.Statuses {
		userTweet := new(UserTweets)
		userTweet.CreatedAt = tweet.CreatedAt
		userTweet.Name = tweet.User.Name
		userTweet.ProfileImage = tweet.User.ProfileImageURL
		userTweet.Text = tweet.Text
		userTweet.Location = tweet.User.Location
		userTweet.ScreenName = tweet.User.ScreenName
		userTweet.Verified = tweet.User.Verified
		userTweet.Retweets = tweet.RetweetCount
		userTweet.Favorites = tweet.FavoriteCount
		userTweet.SentimentRating = model.SentimentAnalysis(tweet.Text, sentiment.English).Score
		tweets = append(tweets, userTweet)
	}

	json.NewEncoder(w).Encode(tweets)
}

// SearchUsers - Searches All Users
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	auth := TwitterAuth()
	httpClient, _ := auth.Configuration()
	client := twitter.NewClient(httpClient)

	vars := mux.Vars(r)
	searchkey, ok := vars["user"]
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
	searchkey, ok := vars["woeid"]

	fmt.Println(searchkey)

	f, err := strconv.ParseInt(searchkey, 10, 64)

	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "The Location of the recent trend couldn't be found}`)
		return
	}

	places := &twitter.TrendsPlaceParams{
		WOEID: f,
	}

	trend, _, err := client.Trends.Place(f, places)
	if err != nil {
		log.Fatal(err)
	}

	trends := []*CurrentTrends{}

	for _, tweet := range trend {

		for _, currentTrend := range tweet.Trends {
			userTrends := new(CurrentTrends)
			userTrends.Name = currentTrend.Name
			userTrends.Volume = currentTrend.TweetVolume
			trends = append(trends, userTrends)
		}
	}

	json.NewEncoder(w).Encode(trends)
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
	searchkey, ok := vars["name"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search}`)
		return
	}

	userTimeline := &twitter.UserTimelineParams{
		// Count:           10,
		ScreenName:      searchkey,
		IncludeRetweets: twitter.Bool(false),
	}

	timeline, _, err := client.Timelines.UserTimeline(userTimeline)
	if err != nil {
		log.Fatal(err)
	}

	tweets := []*UserTimeline{}

	for _, tweet := range timeline {
		userTweet := new(UserTimeline)
		userTweet.ID = tweet.IDStr
		userTweet.Name = tweet.User.Name
		userTweet.ScreenName = tweet.User.ScreenName
		userTweet.Verified = tweet.User.Verified
		userTweet.CreatedAt = tweet.CreatedAt
		userTweet.Location = tweet.User.Location
		userTweet.Text = tweet.Text
		userTweet.TotalTweets = tweet.User.StatusesCount
		userTweet.Retweets = tweet.RetweetCount
		userTweet.Favorites = tweet.FavoriteCount
		userTweet.Followers = tweet.User.FollowersCount
		userTweet.ProfileImage = tweet.User.ProfileImageURL
		userTweet.SentimentRating = model.SentimentAnalysis(tweet.Text, sentiment.English).Score
		tweets = append(tweets, userTweet)
	}

	json.NewEncoder(w).Encode(tweets)
}

// StreamUserTweets streams realtime tweets give keyword
func StreamUserTweets(w http.ResponseWriter, r *http.Request) {
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

	// Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		// json.NewEncoder(w).Encode(tweet.Text)
		fmt.Println(tweet.Text)
		// fmt.Println(tweet.Text)
	}
	fmt.Println("Starting Sream...")

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{searchkey},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
