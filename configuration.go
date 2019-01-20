package main

import (
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

type TwitterAuthentication struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Configuration - Configures Twitter API
func (a *TwitterAuthentication) Configuration() (*http.Client, error) {
	config := oauth1.NewConfig(a.ConsumerKey, a.ConsumerSecret)
	token := oauth1.NewToken(a.AccessToken, a.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return httpClient, nil
}

// TwitterAuth retreives API keys to authorize request from the Twitter API
func TwitterAuth() *TwitterAuthentication {
	return &TwitterAuthentication{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}
}
