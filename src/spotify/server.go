package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/zmb3/spotify"
)

var spotify_client = AuthenticateAndLogUser()

type PlayingData struct {
	Artist   string `json:"artist"`
	Track    string `json:"track"`
	Album    string `json:"album"`
	Progress int    `json:"progress"`
	Length   int    `json:"length"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/nowplaying", NowPlaying)

	log.Fatal(http.ListenAndServe(":8082", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/")
}

// NowPlaying route, responds with a JSON encoded PlayingData object.
func NowPlaying(w http.ResponseWriter, r *http.Request) {
	var current_playing = GetUserCurrentlyPlayingTrack(spotify_client)
	track_data := current_playing.Item
	track := PlayingData{Artist: track_data.Artists[0].Name, Track: track_data.Name, Album: track_data.Album.Name, Progress: current_playing.Progress, Length: track_data.Duration}
	if err := json.NewEncoder(w).Encode(track); err != nil {
		panic(err)
	}
}
