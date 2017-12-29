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
)

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func Truncate(some float64) string {
	val := float64(int(some*100)) / 100
	return FloatToString(val)
}

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	stream := api.PublicStreamFilter(url.Values{
		"follow": []string{"14338147", "928901093974794240", "1656328279", "961445378"},
	})

	defer stream.Stop()

	coinInfo, err := coinApi.GetCoinData("litecoin")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf(Truncate(coinInfo.PriceUsd))
		tweet := "current price: " + (Truncate(coinInfo.PriceUsd)) + " 1hr change: " + (Truncate(coinInfo.PercentChange1h)) + " 24hr change: " + (Truncate(coinInfo.PercentChange24h)) + " #litecoinbot"
		fmt.Printf(tweet)
		api.PostTweet(tweet, nil)

	}

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)

		if !ok {
			logrus.Warningf("recieved unexpected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		// _, err := api.Retweet(t.Id, false)
		// if err != nil {
		// 	logrus.Errorf("cant retweet %d: %v", t.Id, err)
		// 	continue
		// }

		logrus.Infof("retweeted %d", t.Id)
	}
	doEvery(5*time.Second, getLTCprice)
}
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getLTCprice(t time.Time) {
	fmt.Println("hi")
}
