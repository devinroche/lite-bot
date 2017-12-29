package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	coinApi "github.com/miguelmota/go-coinmarketcap"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("TWITTER_ACCESS_SECRET")
	peopleWatch       = []string{"14338147", "928901093974794240", "1656328279", "961445378"}
)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	stream := api.PublicStreamFilter(url.Values{
		"follow": peopleWatch,
	})

	defer stream.Stop()
	go doEvery(10*time.Second, GetLitecoin)

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)

		if !ok {
			logrus.Warningf("recieved unexpected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		_, err := api.Retweet(t.Id, false)
		if err != nil {
			logrus.Errorf("cant retweet %d: %v", t.Id, err)
			continue
		}

		logrus.Infof("retweeted %d", t.Id)
	}
}
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		fmt.Printf("foo")
		f(x)
	}
}

//GetLitecoin data from coinmarketcap
func GetLitecoin(t time.Time) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	coinInfo, err := coinApi.GetCoinData("litecoin")
	if err != nil {
		log.Println(err)
	} else {
		usdprice := strconv.FormatFloat(coinInfo.PriceUsd, 'f', -1, 64)
		percent1hr := strconv.FormatFloat(coinInfo.PercentChange1h, 'f', -1, 64)
		percent24hr := strconv.FormatFloat(coinInfo.PercentChange24h, 'f', -1, 64)

		tweet := "Current Price: $" + usdprice + "\n1 Hour Change: " + percent1hr + "%" + "\n24 Hour Change: " + percent24hr + "%" + "\n#litecoinbot $ltc #litecoin"
		fmt.Printf(tweet)
		api.PostTweet(tweet, nil)
	}
	time.Sleep(1 * time.Second)
}
