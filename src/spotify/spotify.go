package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
)

type Track struct {
	Name    string
	ID      spotify.ID
	AddedAt string
}

type Playlist struct {
	Name string
	ID   spotify.ID
}

type TopArtist string

type UserMusic struct {
	Playlists  []Playlist
	TopArtists []TopArtist
	Tracks     []Track
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

func CurrentUserGetAllPlaylistTracks(client *spotify.Client, playlists []Playlist) []Track {
	var tracks []Track
	for _, playlist := range playlists {
		fmt.Println("Getting tracks for playlist: " + playlist.Name)

		var opt spotify.Options
		opt.Limit = new(int)
		*opt.Limit = 50 // Maximum amount possible in one request.
		opt.Offset = new(int)
		*opt.Offset = 0

		playlist_track_page, err := client.GetPlaylistTracksOpt(playlist.ID, &opt, "items(added_at,track(name,id))")
		for _, playlist_track := range playlist_track_page.Tracks {
			tracks = append(tracks, Track{playlist_track.Track.Name, playlist_track.Track.ID, playlist_track.AddedAt})
		}

		if err != nil {
			fmt.Printf("Error retrieving playlist tracks: %s\n", err)
		} else {
			for ; *opt.Offset < playlist_track_page.Total; *opt.Offset += len(playlist_track_page.Tracks) {
				playlist_track_page, err = client.GetPlaylistTracksOpt(playlist.ID, &opt, "items(added_at,track(name,id))")

				for _, playlist_track := range playlist_track_page.Tracks {
					tracks = append(tracks, Track{playlist_track.Track.Name, playlist_track.Track.ID, playlist_track.AddedAt})
				}
			}
		}
	}
	return tracks
}

func CurrentUserGetAllPlaylists(client *spotify.Client) []Playlist {
	// Retrieve all of the users playlists
	var opt spotify.Options
	opt.Limit = new(int)
	*opt.Limit = 50 // Maximum amount possible in one request.
	opt.Offset = new(int)
	*opt.Offset = 0

	var playlists []Playlist

	playlist_page, err := client.CurrentUsersPlaylistsOpt(&opt)
	for _, playlist := range playlist_page.Playlists {
		playlists = append(playlists, Playlist{playlist.Name, playlist.ID})
	}
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
	} else {
		for _, artist := range top_artist_page.Artists {
			top_artists = append(top_artists, TopArtist(artist.Name))
		}
	}

	return top_artists
}

func GetAllCurrentUserMusic(client *spotify.Client) UserMusic {
	var user_music UserMusic
	user_music.Playlists = CurrentUserGetAllPlaylists(client)
	user_music.TopArtists = CurrentUserGetTopArtists(client)
	user_music.Tracks = CurrentUserGetAllPlaylistTracks(client, user_music.Playlists)
	return user_music
}

func GetUserCurrentlyPlayingTrack(client *spotify.Client) *spotify.CurrentlyPlaying {
	currently_playing, err := client.PlayerCurrentlyPlaying()

	if err != nil {
		fmt.Printf("Error retrieving currently played song: %s\n", err)
		return nil
	} else {
		return currently_playing
	}
}
