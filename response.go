package main

type UserTweets struct {
	CreatedAt       string
	Name            string
	ProfileImage    string
	ScreenName      string
	Location        string
	Verified        bool
	Text            string
	Retweets        int
	Favorites       int
	SentimentRating uint8
}

type UserTimeline struct {
	ID              string
	Name            string
	ScreenName      string
	Verified        bool
	ProfileImage    string
	CreatedAt       string
	Location        string
	Text            string
	TotalTweets     int
	Followers       int
	Retweets        int
	Favorites       int
	SentimentRating uint8
}

type CurrentTrends struct {
	Name   string
	Volume int64
}
