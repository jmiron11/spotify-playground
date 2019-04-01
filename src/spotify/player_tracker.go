package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

const TRACK_INTERVAL = 20 //seconds

func PlayerTrack(client *spotify.Client) {
	currently_playing := GetUserCurrentlyPlayingTrack(client)
	if currently_playing == nil || !currently_playing.Playing {
		fmt.Println("Did not detect player has track currently playing")
	} else {
		track := currently_playing.Item
		fmt.Printf("User is playing the song %s, t: %d, p: %d\n", track.Name, currently_playing.Timestamp, currently_playing.Progress)
	}
}

func AuthenticateAndLogUser() *spotify.Client {
	client := Authenticate()
	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	return client
}

func StartIntervalPlayerTracker(client *spotify.Client) {
	ticker := time.NewTicker(TRACK_INTERVAL * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				PlayerTrack(client)
			case <-quit:
				ticker.Stop()
				fmt.Println("Ticker stopped")
				return
			}
		}
	}()
}

func InitializeMongodb() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	return session
}

func TeardownMongodb(session *mgo.Session) {
	if session != nil {
		session.Close()
	}
}

type Person struct {
	Name  string
	Phone string
}

func main() {
	// spotify_client := AuthenticateAndLogUser()
	// StartIntervalPlayerTracker(spotify_client)
	session := InitializeMongodb()

	c := session.DB("test").C("people")
	err := c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	TeardownMongodb(session)
}
