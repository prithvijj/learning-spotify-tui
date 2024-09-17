package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

var (
	auth  = spotify.NewAuthenticator("http://localhost:8080/callback", spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopePlaylistReadPrivate)
	state = "abc123"
	ch    = make(chan *spotify.Client)
)

func spotifyAuthentication(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "could not get spotify oauth token", http.StatusForbidden)
		log.Fatal(err)
	}

	spotifyClient := auth.NewClient(token)
	fmt.Fprintf(w, "Login Completed")

	ch <- &spotifyClient

}

func getUserPlayerState(client *spotify.Client) (*spotify.PlayerState, error) {
	playerState, err := client.PlayerState()
	if err != nil {
		return nil, err
	}

	return playerState, nil

}

func loadTracksForPlaylist(client *spotify.Client, playlistID spotify.ID) ([]spotify.FullTrack, error) {
	playlistTracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return nil, err
	}

	var tracks []spotify.FullTrack
	for _, item := range playlistTracks.Tracks {
		track := item.Track
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func initializePlaylists(client *spotify.Client) ([]spotify.SimplePlaylist, error) {
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return nil, err
	}

	return playlists.Playlists, nil
}
