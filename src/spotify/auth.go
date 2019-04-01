/*
 * auth.go - Justin Miron
 * This wraps the spotify authentication behavior, to retrieve
 * an authenticated client which cay be used with the Spotify
 * web api.
 */
package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

// This redirectURI must be registered with Spotify, matching the
// CLIENT_KEY, CLIENT_SECRET specified in the environmental variables
// The uri must also match the HTTP server setup later.
const redirectURI = "http://localhost:8080/authcallback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserTopRead, spotify.ScopePlaylistReadPrivate, spotify.ScopeUserReadPlaybackState)
	cb_ch = make(chan *spotify.Client)
	state = "completelyrandomstatevalue"
)

// Authentication with Spotify creates a http server for listening to
// the authentication callback.
func Authenticate() *spotify.Client {
	// Start an HTTP server for listening to the callback URI
	http.HandleFunc("/authcallback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	// Get the spotify authorization URL and pass it to the user, wait for authorization
	// callback via waiting on channel (cb_ch).
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-cb_ch
	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	cb_ch <- &client
}
