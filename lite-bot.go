package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("TWITTER_ACCESS_SECRET")
)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := &logger{logrus.New()}
	api.SetLogger(log)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"litecoin"},
	})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)

		if !ok {
			log.Warningf("recieved unexpected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		_, err := api.Retweet(t.Id, false)
		if err != nil {
			log.Errorf("cant retweet %d: %v", t.Id, err)
			continue
		}

		log.Infof("retweeted %d", t.Id)
	}
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{}) {
	log.Error(args...)
}
func (log *logger) Criticalf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
func (log *logger) Notice(args ...interface{}) {
	log.Info(args...)
}
func (log *logger) Noticef(format string, args ...interface{}) {
	log.Infof(format, args...)
}
