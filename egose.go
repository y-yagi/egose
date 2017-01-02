package main

import (
	"flag"
	"fmt"
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

func main() {
	config, err := loadConfig()

	if err != nil {
		fmt.Printf("Config file load Error: %v\nPlease create a config file.\n", err)
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

	oauthConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
		Count: count,
	})

	if err != nil {
		fmt.Printf("Search API Error:%v\n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User", "Text", "URL"})
	for _, tweet := range search.Statuses {
		table.Append([]string{tweet.User.Name, runewidth.Truncate(tweet.Text, 80, "..."), tweetURL(tweet)})
	}
	table.Render()
}
