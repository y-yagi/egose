package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func tweetURL(tweet twitter.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.ID)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	var query string
	var count int
	flag.StringVar(&query, "q", "", "Search query")
	flag.IntVar(&count, "c", 10, "Search count")
	flag.Parse()

	if len(query) == 0 {
		fmt.Println("Please specify search query.")
		os.Exit(1)
	}

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
		Count: count,
	})

	if err != nil {
		fmt.Printf("Search API Error:%v\n", err)
		os.Exit(1)
	}

	for _, tweet := range search.Statuses {
		fmt.Printf("%v: %v( %v )\n", tweet.User.Name, tweet.Text, tweetURL(tweet))
	}
}
