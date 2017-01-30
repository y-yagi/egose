package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/olekukonko/tablewriter"
)

type config struct {
	TwitterConsumerKey    string `yaml:"twitterConsumerKey"`
	TwitterConsumerSecret string `yaml:"twitterConsumerSecret"`
	TwitterAccessToken    string `yaml:"twitterAccessToken"`
	TwitterAccessSecret   string `yaml:"twitterAccessSecret"`
}

func loadConfig() (*config, error) {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")
	}

	fname := filepath.Join(home, ".config", "egose", "config.yml")
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	var cfg config
	err = yaml.Unmarshal(buf, &cfg)
	return &cfg, err
}

func tweetURL(tweet twitter.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.ID)
}

func buildTwitterClient(cfg *config) *twitter.Client {
	oauthConfig := oauth1.NewConfig(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func getTimelineTweets(client *twitter.Client, count int) ([]twitter.Tweet, error) {
	tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})
	return tweets, err
}

func getUserTimelineTweets(client *twitter.Client, screenName string, count int) ([]twitter.Tweet, error) {
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: screenName,
		Count:      count,
	})
	return tweets, err
}

func searchTweets(client *twitter.Client, count int, query string) ([]twitter.Tweet, error) {
	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
		Count: count,
	})
	if err != nil {
		return nil, err
	}
	return search.Statuses, nil
}

func updateStatus(client *twitter.Client, status string) error {
	_, _, err := client.Statuses.Update(status, nil)
	return err
}

func main() {
	config, err := loadConfig()

	if err != nil {
		fmt.Printf("Config file load Error: %v\nPlease create a config file.\n", err)
		os.Exit(1)
	}

	var query string
	var user string
	var count int
	var status string
	var tweets []twitter.Tweet

	flag.StringVar(&query, "q", "", "Search query")
	flag.StringVar(&user, "u", "", "Show user timeline")
	flag.IntVar(&count, "c", 50, "Search count")
	flag.StringVar(&status, "p", "", "Post tweet")
	flag.Parse()

	client := buildTwitterClient(config)

	if len(status) > 0 {
		err = updateStatus(client, status)
		if err != nil {
			fmt.Printf("twitter API Error:%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(query) > 0 {
		tweets, err = searchTweets(client, count, query)
	} else if len(user) > 0 {
		tweets, err = getUserTimelineTweets(client, user, count)
	} else {
		tweets, err = getTimelineTweets(client, count)
	}

	if err != nil {
		fmt.Printf("twitter API Error:%v\n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User", "Text", "URL"})
	for _, tweet := range tweets {
		table.Append([]string{tweet.User.Name, runewidth.Truncate(html.UnescapeString(tweet.Text), 80, "..."), tweetURL(tweet)})
	}
	table.Render()
}
