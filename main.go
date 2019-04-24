package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dghubble/go-twitter/twitter"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/y-yagi/configure"
	"github.com/y-yagi/gocui"
)

var tweets []twitter.Tweet

const cmd = "egose"

// Config manage config info
type Config struct {
	TwitterConsumerKey    string `toml:"twitterConsumerKey"`
	TwitterConsumerSecret string `toml:"twitterConsumerSecret"`
	TwitterAccessToken    string `toml:"twitterAccessToken"`
	TwitterAccessSecret   string `toml:"twitterAccessSecret"`
}

func tweetURL(tweet twitter.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.ID)
}

func readTweetFromFile() (string, error) {
	const defaultEditor = "vi"

	msgFile := filepath.Join(configure.ConfigDir(cmd), "TWEET")

	// Clean up file
	_ = os.Remove(msgFile)
	cmd := exec.Command(defaultEditor, msgFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	dat, _ := ioutil.ReadFile(msgFile)
	return string(dat), nil
}

func updateStatus(egose *Egose) error {
	var tweet string
	var err error

	if len(flag.Args()) > 0 {
		tweet = flag.Args()[0]
	} else {
		tweet, err = readTweetFromFile()
		if err != nil {
			return errors.Wrap(err, "unexpected error")
		}
	}

	if len(tweet) == 0 {
		return nil
	}

	err = egose.UpdateStatus(tweet)
	if err != nil {
		msg := fmt.Sprintf("twitter API Error.\ntweet: %v", tweet)
		return errors.Wrap(err, msg)
	}

	return nil
}

func showTweetsWithTable(tweets []twitter.Tweet) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User", "Text", "URL"})
	for _, tweet := range tweets {
		table.Append([]string{runewidth.Truncate(tweet.User.Name, 30, "..."), runewidth.Truncate(html.UnescapeString(tweet.Text), 80, "..."), tweetURL(tweet)})
	}
	table.Render()
}

func showTweetsWithGui() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return errors.Wrap(err, "gui create error")
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		return errors.Wrap(err, "key bindings error")
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return errors.Wrap(err, "unexpected error")
	}
	return nil
}

func main() {
	var cfg Config

	err := configure.Load(cmd, &cfg)
	if err != nil {
		fmt.Printf("Config file load Error: %v\nPlease create a config file.\n", err)
		os.Exit(1)
	}

	var query string
	var user string
	var count int
	var status bool

	flag.StringVar(&query, "q", "", "Search query")
	flag.StringVar(&user, "u", "", "Show user timeline")
	flag.IntVar(&count, "c", 50, "Search count")
	flag.BoolVar(&status, "p", false, "Post tweet. If you specify a message, that message will be sent as is. If you do not specify a message, the editor starts up.")
	flag.Parse()

	egose := NewEgose(&cfg)

	if status {
		err = updateStatus(egose)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(query) > 0 {
		tweets, err = egose.SearchTweets(count, query)
	} else if len(user) > 0 {
		tweets, err = egose.GetUserTimelineTweets(user, count)
	} else {
		tweets, err = egose.GetTimelineTweets(count)
	}

	if err != nil {
		fmt.Printf("twitter API Error:%v\n", err)
		os.Exit(1)
	}

	if err = showTweetsWithGui(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
