package main

import "time"

// TwitchStreamOnline is a full representation of an online twitch stream json object
type TwitchStreamOnline struct {
	Links struct {
		Channel string `json:"channel"`
		Self    string `json:"self"`
	} `json:"_links"`
	Stream struct {
		Game        string    `json:"game"`
		Viewers     int       `json:"viewers"`
		AverageFps  float64   `json:"average_fps"`
		Delay       int       `json:"delay"`
		VideoHeight int       `json:"video_height"`
		IsPlaylist  bool      `json:"is_playlist"`
		CreatedAt   time.Time `json:"created_at"`
		ID          int64     `json:"_id"`
		Channel     struct {
			Mature                       bool        `json:"mature"`
			Status                       string      `json:"status"`
			BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName                  string      `json:"display_name"`
			Game                         string      `json:"game"`
			Delay                        interface{} `json:"delay"`
			Language                     string      `json:"language"`
			ID                           int         `json:"_id"`
			Name                         string      `json:"name"`
			CreatedAt                    time.Time   `json:"created_at"`
			UpdatedAt                    time.Time   `json:"updated_at"`
			Logo                         string      `json:"logo"`
			Banner                       string      `json:"banner"`
			VideoBanner                  string      `json:"video_banner"`
			Background                   interface{} `json:"background"`
			ProfileBanner                string      `json:"profile_banner"`
			ProfileBannerBackgroundColor string      `json:"profile_banner_background_color"`
			Partner                      bool        `json:"partner"`
			URL                          string      `json:"url"`
			Views                        int         `json:"views"`
			Followers                    int         `json:"followers"`
			Links                        struct {
				Self          string `json:"self"`
				Follows       string `json:"follows"`
				Commercial    string `json:"commercial"`
				StreamKey     string `json:"stream_key"`
				Chat          string `json:"chat"`
				Features      string `json:"features"`
				Subscriptions string `json:"subscriptions"`
				Editors       string `json:"editors"`
				Teams         string `json:"teams"`
				Videos        string `json:"videos"`
			} `json:"_links"`
		} `json:"channel"`
		Preview struct {
			Small    string `json:"small"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Template string `json:"template"`
		} `json:"preview"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"stream"`
}

// TwitchStreamOffline is a full representation of an offline twitch stream json object
/*type TwitchStreamOffline struct {
	Stream interface{} `json:"stream"`
	Links  struct {
		Self    string `json:"self"`
		Channel string `json:"channel"`
	} `json:"_links"`
}
*/

// TwitchRoot is a struct to hold a json reponse from a call to twitch's root api
type TwitchRoot struct {
	Identified bool `json:"identified"`
	Links      struct {
		User    string `json:"user"`
		Channel string `json:"channel"`
		Search  string `json:"search"`
		Streams string `json:"streams"`
		Ingests string `json:"ingests"`
		Teams   string `json:"teams"`
	} `json:"_links"`
	Token struct {
		Valid         bool        `json:"valid"`
		Authorization interface{} `json:"authorization"`
	} `json:"token"`
}

// TwitchStream is a generic tester
type TwitchStream struct {
	Stream interface{} `json:"stream"`
	Links  interface{} `json:"_links"`
}
