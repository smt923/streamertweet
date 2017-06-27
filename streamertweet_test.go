package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func TestReadStreamersFunc(t *testing.T) {
	t.Log("\tTesting if our ReadStreamers function works correctly")
	_, err := ReadStreamers(streamersFile)
	if err != nil {
		t.Error("Error was not nil for reading streamers: ", err)
	}
}

func TestUsersTweetBool(t *testing.T) {
	t.Log("\tTesting if we have our streamers formatted and set up correctly")
	streams, err := ReadStreamers(streamersFile)
	if err != nil {
		t.Error("Error was not nil for reading streamers: ", err)
	}
	for i, streamer := range streams {
		if streamer.Tweeted != true {
			t.Error("Streamer not marked as tweeted - we would be tweeting on first run!", streams[i])
		}
		if streamer.Name == "" {
			t.Error("Streamer was missing Name", streams[i])
		}
		if streamer.Twitter == "" {
			t.Error("Streamer was missing Twitter", streams[i])
		}
		if streamer.User == "" {
			t.Error("Streamer was missing User", streams[i])
		}
	}
}

// Test if our twitter creds work and we can auth with twitter
func TestTwitter(t *testing.T) {
	keys, err := godotenv.Read()
	if err != nil {
		t.Error(err)
	}
	t.Log("\tTesting if we can our authentication tokens are valid and working with the Twitter.com API")

	anaconda.SetConsumerKey(keys["TWITKEY"])
	anaconda.SetConsumerSecret(keys["TWITSEC"])
	client := anaconda.NewTwitterApi(keys["TOK"], keys["TOKSEC"])

	ok, err := client.VerifyCredentials()
	if err != nil {
		t.Log(ok)
		t.Error(err)
	}
	if ok != true {
		t.Log(ok)
		t.Error(err)
	}

}

func TestTwitch(t *testing.T) {
	keys, err := godotenv.Read()
	if err != nil {
		t.Error(err)
	}

	t.Log("\tTesting if we can authenticate our Client-ID with Twitch.tv API")
	twitchtestjson := new(TwitchRoot)
	twitchtest := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken"+"/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Accept", twitchmime)
	req.Header.Add("Client-ID", keys["TWITCHID"])

	r, err := twitchtest.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(twitchtestjson)

	if twitchtestjson.Identified != true {
		t.Error("Twitch could not authenticate us")
	}
}

// test our length check functions as we expect
func TestLengthCheck(t *testing.T) {
	t.Log("\tTesting if we can check a tweet's length correctly")
	longString := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis p"
	shortString := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor."
	exactString := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et ma"

	if FitsTweet(longString) != false {
		t.Error("FitsTweet allowed a long string")
	}
	if FitsTweet(shortString) != true {
		t.Error("FitsTweet didn't allow a string that fits")
	}
	if FitsTweet(exactString) != true {
		t.Error("FitsTweet didn't allow a 140 chars string")
	}
}

// test we can shorten games as we expect
func TestShortGame(t *testing.T) {
	t.Log("\tTesting if we can shorten a game's name correctly")
	if ShortGame("Counter-Strike: Global Offensive") != "#CSGO" {
		t.Error("ShortGame did not return a shortened string for CSGO")
	}
	if ShortGame("Hearthstone: Heroes of Warcraft") != "#Hearthstone" {
		t.Error("ShortGame did not return a shortened string for Hearthstone")
	}
	if ShortGame("Overwatch") != "#Overwatch" {
		t.Error("ShortGame did not return a hashtag version for Overwatch")
	}
	if ShortGame("Dark Souls") != "Dark Souls" {
		t.Error("ShortGame did not simply return the game name if it's not in our list")
	}
	if ShortGame("") != "some games" {
		t.Error("ShortGame did not return 'some games' for an empty game")
	}
}
