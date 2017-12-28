package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

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

	stream := api.PublicStreamFilter(url.Values{
		"follow": []string{"14338147", "928901093974794240", "1656328279", "961445378"},
	})

	defer stream.Stop()
	response, _ := http.Get("https://api.coinmarketcap.com/v1/ticker/litecoin/")
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

	api.PostTweet("foo", nil)
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
	// fmt.Printf("%v: Hello, World!\n", t)
	response, _ := http.Get("https://api.coinmarketcap.com/v1/ticker/litecoin/")
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
}
