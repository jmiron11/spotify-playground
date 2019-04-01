package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
)

type Playlist struct {
	Name string
	ID   spotify.ID
}

type TopArtist string

type UserMusic struct {
	Playlists  []Playlist
	TopArtists []TopArtist
}

func CurrentUserGetAllPlaylists(client *spotify.Client) []Playlist {
	// Retrieve all of the users playlists
	var opt spotify.Options
	opt.Limit = new(int)
	*opt.Limit = 50
	opt.Offset = new(int)
	*opt.Offset = 0

	var playlists []Playlist

	playlist_page, err := client.CurrentUsersPlaylistsOpt(&opt)
	if err != nil {
		fmt.Printf("Error retrieving user playlists: %s\n", err)
	} else {
		for ; *opt.Offset < playlist_page.Total; *opt.Offset += len(playlist_page.Playlists) {
			playlist_page, err = client.CurrentUsersPlaylistsOpt(&opt)

			for _, playlist := range playlist_page.Playlists {
				playlists = append(playlists, Playlist{playlist.Name, playlist.ID})
			}
		}
	}
	return playlists
}

func CurrentUserGetTopArtists(client *spotify.Client) []TopArtist {
	var top_artists []TopArtist
	top_artist_page, err := client.CurrentUsersTopArtists()
	if err != nil {
		fmt.Printf("Error retrieving top artists: %s\n", err)
		fmt.Print(err)
	} else {
		for _, artist := range top_artist_page.Artists {
			top_artists = append(top_artists, TopArtist(artist.Name))
		}
	}

	return top_artists
}

func main() {
	client := Authenticate()
	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	var user_music UserMusic
	user_music.Playlists = CurrentUserGetAllPlaylists(client)
	user_music.TopArtists = CurrentUserGetTopArtists(client)
}
