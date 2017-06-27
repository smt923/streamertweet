# Streamertweet

A bot that will check if streamers are online, and if so, tweet about them streaming

## Usage

Simply build and run, the streamers.txt (see below) by default should be in the working directory, the bot will automatically begin checking periodically

## Streamers.txt format

Streamers are read by default from streamers.txt in the main directory, this is just a simple csv file with the format:
```
general alias, twitch username, @twitter handle
```
Each streamer should be seperated by a new line

## .env / environment variables

The envionrment variables needed are as follows:
```
# twitch
TWITCHID=Twitch api client ID
	
# twitter
TWITKEY=Twitter consumer key
TWITSEC=Twitter consumer secret
TOK=Twitter api token
TOKSEC=Twitter api token secret
```