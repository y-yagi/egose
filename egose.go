package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Egose struct {
	client *twitter.Client
}

func NewEgose(cfg *Config) *Egose {
	var egose Egose
	egose.init(cfg)
	return &egose
}

func (e *Egose) init(cfg *Config) {
	e.client = buildTwitterClient(cfg)
}

func buildTwitterClient(cfg *Config) *twitter.Client {
	oauthConfig := oauth1.NewConfig(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func (e *Egose) GetTimelineTweets(count int) ([]twitter.Tweet, error) {
	tweets, _, err := e.client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})
	return tweets, err
}

func (e *Egose) GetUserTimelineTweets(screenName string, count int) ([]twitter.Tweet, error) {
	tweets, _, err := e.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: screenName,
		Count:      count,
	})
	return tweets, err
}

func (e *Egose) SearchTweets(count int, query string) ([]twitter.Tweet, error) {
	search, _, err := e.client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
		Count: count,
	})
	if err != nil {
		return nil, err
	}
	return search.Statuses, nil
}

func (e *Egose) UpdateStatus(status string) error {
	_, _, err := e.client.Statuses.Update(status, nil)
	return err
}
