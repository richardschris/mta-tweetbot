package main

import (
	"bufio"
	"fmt"
	"html"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reasons := make([]string, 0)
	for scanner.Scan() {
		reasons = append(reasons, scanner.Text())
	}
	return reasons
}

func mtaTweetListener(client *twitter.Client) {
	becauseReasons := readFile("because.txt")
	resumedReasons := readFile("resumed.txt")
	rand.Seed(time.Now().Unix())

	params := &twitter.UserTimelineParams{
		ScreenName:     "NYCTSubway",
		Count:          10,
		TweetMode:      "extended",
		ExcludeReplies: twitter.Bool(true),
	}

	for {
		tweets, _, _ := client.Timelines.UserTimeline(params)
		fmt.Printf("Found %v Tweets\n", len(tweets))
		sinceIDs := make([]int64, 0)
		for _, tweet := range tweets {
			if strings.Contains(tweet.FullText, "because") {
				createNewTweet(tweet.FullText, becauseReasons, "because", client)
			} else if strings.Contains(tweet.FullText, "resumed") {
				createNewTweet(tweet.FullText, resumedReasons, "resumed", client)
			}
			sinceIDs = append(sinceIDs, tweet.ID)
		}
		if len(tweets) > 0 {
			finalSinceID := sinceIDs[0]
			for _, sinceID := range sinceIDs {
				if finalSinceID < sinceID {
					finalSinceID = sinceID
				}
			}
			params.SinceID = finalSinceID
		}
		time.Sleep(60 * time.Second)
		fmt.Printf("Begin again with SinceID %v \n", params.SinceID)
	}
}

func createNewTweet(reasonString string, reasons []string, reasonType string, client *twitter.Client) {
	reasonIndex := strings.Index(reasonString, reasonType)
	tweetSlice := reasonString[0 : reasonIndex+len(reasonType)]
	reasonTweet := tweetSlice + reasons[rand.Intn(len(reasons))]
	client.Statuses.Update(html.UnescapeString(reasonTweet), nil)
	fmt.Printf("Made a tweet at %v\n", time.Now())
}

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// boilerplate for go-twitter
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		fmt.Println("WOAH ERROR WITH CREDENTIALS")
		panic(err)
	}
	mtaTweetListener(client)
}
