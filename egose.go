package main

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Egose manage twitter related processing
type Egose struct {
	client *twitter.Client
}

// NewEgose generate new Egose
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

// GetTimelineTweets get timeline tweet
func (e *Egose) GetTimelineTweets(count int) ([]twitter.Tweet, error) {
	tweets, _, err := e.client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})
	return tweets, err
}

func (e *Egose) GetDebugTweets(count int) ([]twitter.Tweet, error) {
	var tweets []twitter.Tweet
	for i := 0; i < count; i++ {
		u := twitter.User{Name: strconv.Itoa(i)}
		tweet := twitter.Tweet{Text: "text" + strconv.Itoa(i), User: &u}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

// GetUserTimelineTweets get specified user timeline tweets
func (e *Egose) GetUserTimelineTweets(screenName string, count int) ([]twitter.Tweet, error) {
	tweets, _, err := e.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: screenName,
		Count:      count,
	})
	return tweets, err
}

// SearchTweets search tweets
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

// UpdateStatus update my curent status
func (e *Egose) UpdateStatus(status string) error {
	_, _, err := e.client.Statuses.Update(status, nil)
	return err
}

// GetListTweets get list Statuses
func (e *Egose) GetListTweets(listID string, count int) ([]twitter.Tweet, error) {
	id, _ := strconv.ParseInt(listID, 10, 64)
	tweets, _, err := e.client.Lists.Statuses(&twitter.ListsStatusesParams{
		ListID: id,
		Count:  count,
	})
	return tweets, err
}

func (e *Egose) GetListMembers(listID string) (*twitter.Members, error) {
	id, _ := strconv.ParseInt(listID, 10, 64)
	members, _, err := e.client.Lists.Members(&twitter.ListsMembersParams{
		ListID: id,
		Count:  100, // TODO: allow to specify via command line
	})
	return members, err
}
