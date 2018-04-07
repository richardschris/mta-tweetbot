package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func mtaTweetListener(client *twitter.Client) {
	reasons := make([]string, 0)
	reasons = append(reasons,
		"nobody knows how anything works anymore",
		"Andrew Cuomo wants you to have a bad day",
		"the universe is conspiring against you",
		"the MTA thinks you can fucking walk",
		"it's infrastructure week",
		"you should eat at Arby's",
		"we replaced all trains with golf carts",
		"we can't have nice things",
		"you should stop and smell the flowers",
		"",
	)

	params := &twitter.UserTimelineParams{
		ScreenName:     "NYCTSubway",
		Count:          10,
		TweetMode:      "extended",
		ExcludeReplies: twitter.Bool(true),
	}
	client.Statuses.Update("TESTING A TWEET", nil)
	for {
		tweets, _, _ := client.Timelines.UserTimeline(params)
		fmt.Printf("Found %v Tweets\n", len(tweets))
		for _, tweet := range tweets {
			if strings.Contains(tweet.FullText, "because") || true {
				createNewTweet(tweet.FullText, reasons, client)
			}
		}
		if len(tweets) > 0 {
			sinceID := tweets[len(tweets)-1]
			params.SinceID = sinceID.ID
		}
		time.Sleep(60 * time.Second)
		fmt.Println("Begin again")
	}
}

func createNewTweet(becauseString string, reasons []string, client *twitter.Client) {
	rand.Seed(time.Now().Unix())
	becauseIndex := strings.Index(becauseString, "because")
	tweetSlice := becauseString[0:becauseIndex]
	becauseTweet := tweetSlice + reasons[rand.Intn(len(reasons))]
	client.Statuses.Update(becauseTweet, nil)
	fmt.Println("Made a tweet at %v", time.Now())
}

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests

	httpClient := config.Client(oauth1.NoContext, token) // Twitter client
	client := twitter.NewClient(httpClient)
	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	client.Accounts.VerifyCredentials(verifyParams)
	mtaTweetListener(client)
}
