package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

const (
	streamersFile = "streamers.txt"
	twitchmime    = `application/vnd.twitchtv.v3+json`
)

// Streamer is a basic object for holding a team member's information
type Streamer struct {
	Name, User, Twitter string
	Tweeted             bool
}

func main() {
	keys, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(keys["TWITKEY"])
	anaconda.SetConsumerSecret(keys["TWITSEC"])
	client := anaconda.NewTwitterApi(keys["TOK"], keys["TOKSEC"])

	streams, err := ReadStreamers(streamersFile)
	if err != nil {
		log.Println(err)
	}

	for {
		fmt.Println("[LOG] - Beginning checks")
		for i, streamer := range streams {
			rand.Seed(time.Now().UnixNano())
			status := checkOnline(streamer.User, keys)
			log.Println("[LOG] - Checked: " + streamer.User + " - " + status + " - Tweeted: " + strconv.FormatBool(streamer.Tweeted))
			if status == "online" && streamer.Tweeted == false {
				streamobj, err := getStream(streamer.User, keys)
				if err != nil {
					continue
				}
				// If the streamer has no twitter, we'll fill it in with their name
				if streamer.Twitter == "NoTwitter" || streamer.Twitter == "" || streamer.Twitter == "none" {
					streamer.Twitter = streamer.Name
				}
				game := ShortGame(streamobj.Stream.Game)
				var tweets = []string{
					fmt.Sprintf("Looking for streams? %s is now live at https://twitch.tv/%s - come and watch!", streamer.Twitter, streamer.User),
					fmt.Sprintf("Come and watch %s with us over on https://twitch.tv/%s ! ", streamer.Twitter, streamer.User),
					fmt.Sprintf("Searching for new streams? %s is now live at https://twitch.tv/%s - come say hi!", streamer.Twitter, streamer.User),
					fmt.Sprintf("Watch %s with us over on https://twitch.tv/%s ! ", streamer.Twitter, streamer.User),
					fmt.Sprintf("Our friend %s is now live at https://twitch.tv/%s - check out the stream!", streamer.Twitter, streamer.User),
					fmt.Sprintf("Check out %s live with %s now! https://twitch.tv/%s ", streamer.Twitter, game, streamer.User),
					fmt.Sprintf("%s and chill with %s ! https://twitch.tv/%s ", game, streamer.Twitter, streamer.User),
					fmt.Sprintf("Want a new streamer to watch? %s is now live at https://twitch.tv/%s - come hang out!", streamer.Twitter, streamer.User),
					fmt.Sprintf("Watch %s play %s over on https://twitch.tv/%s ", streamer.Twitter, game, streamer.User),
				}
				fallbacktweet := fmt.Sprintf("Come and watch %s with us over on https://twitch.tv/%s !", streamer.Twitter, streamer.User)
				randtweet := tweets[rand.Intn(len(tweets))]

				// first check if [test] is in the stream title, if not, continue
				if !strings.Contains(streamobj.Stream.Channel.Status, "[test]") {
					if FitsTweet(randtweet) {
						tweet, err := client.PostTweet(randtweet, nil)
						if err != nil {
							log.Println(err)
						}
						fmt.Println("Tweet sent: ", tweet.Text)
						streams[i].Tweeted = true
					} else {
						tweet, err := client.PostTweet(fallbacktweet, nil)
						if err != nil {
							log.Println(err)
						}
						fmt.Println("Tweet did not fit, sent: ", tweet.Text)
						streams[i].Tweeted = true
					}
				}

			} else if status == "offline" && streamer.Tweeted == true {
				log.Println("\t[LOG] - Online streamer gone offline, or I just started up: " + streamer.Name)
				streams[i].Tweeted = false
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("[LOG] - Waiting for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}

func checkOnline(channel string, keys map[string]string) string {
	stream := new(TwitchStream)
	twitchclient := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/streams/"+channel, nil)
	if err != nil {
		log.Println(err)
		return "offline"
	}

	req.Header.Add("Accept", twitchmime)
	req.Header.Add("Client-ID", keys["TWITCHID"])

	r, err := twitchclient.Do(req)
	if err != nil {
		log.Println(err)
		return "offline"
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(stream)
	switch stream.Stream {
	case nil:
		return "offline"
	default:
		return "online"
	}
}

func getStream(channel string, keys map[string]string) (*TwitchStreamOnline, error) {
	stream := new(TwitchStreamOnline)
	twitchclient := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/streams/"+channel, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Add("Accept", twitchmime)
	req.Header.Add("Client-ID", keys["TWITCHID"])

	r, err := twitchclient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(stream)
	return stream, nil
}

// FitsTweet checks if an input string would fit into a tweet
func FitsTweet(tweet string) bool {
	return len(tweet) <= 140
}

// ShortGame shorterns the length of common games on Twitch.tv, adding a hashtag
func ShortGame(game string) string {
	switch game {
	case "Hearthstone: Heroes of Warcraft":
		return "#Hearthstone"
	case "Counter-Strike: Global Offensive":
		return "#CSGO"
	case "Return to Castle Wolfenstein":
		return "#RTCW"
	case "Call of Duty 4: Modern Warfare":
		return "#CoD4"
	case "Team Fortress 2":
		return "#TF2"
	case "World of Warcraft":
		return "#WoW"
	case "Overwatch":
		return "#Overwatch"
	case "Quake Live":
		return "#QuakeLive"
	case "":
		return "some games"
	default:
		return game
	}
}

// ReadStreamers will read a given csv file and extract streamer information from it in this format:
// Name, Twitch user, Twitter
func ReadStreamers(file string) ([]Streamer, error) {
	var streamError error

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	reader.Comment = '#'

	result, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	streams := make([]Streamer, len(result))

	if err == io.EOF {
		return streams, streamError
	}
	// default our streamers to true so we dont tweet them out when we first run
	for i, _streamer := range result {
		streams[i].Name = _streamer[0]
		streams[i].User = _streamer[1]
		streams[i].Twitter = _streamer[2]
		streams[i].Tweeted = true
	}
	return streams, streamError
}
