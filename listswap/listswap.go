package listswap

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// Run runs the tool to create a new playlist
func Run(playlistID string) {

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	playlistSID := spotify.ID(playlistID)
	ctx := context.Background()

	playList, err := client.GetPlaylist(
		ctx,
		playlistSID,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Making a new list based on \"%v\" with %v tracks", playList.Name, playList.Tracks.Total)

	newPlaylist, err := client.CreatePlaylistForUser(ctx, playList.Owner.ID, playList.Name+" (2)", playList.Name+" - with swaps", false, false)
	if err != nil {
		log.Fatal(err)
	}

	tracks, err := client.GetPlaylistItems(
		ctx,
		playlistSID,
	)
	if err != nil {
		log.Fatal(err)
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for page := 1; ; page++ {
		newTracks := make([]spotify.ID, 0, len(tracks.Items))
		for _, track := range tracks.Items {
			t := track.Track.Track

			fullAlbum, err := client.GetAlbum(
				ctx,
				spotify.ID(t.Album.ID),
			)
			if err != nil {
				log.Fatal(err)
			}

			trackCount := len(fullAlbum.Tracks.Tracks)
			newIdx := r.Intn(trackCount)
			newT := fullAlbum.Tracks.Tracks[newIdx]

			log.Printf("  Swapping %v from %v for %v", t.Name, t.Album.Name, newT.Name, newT.ID)
			newTracks = append(newTracks, newT.ID)
		}
		_, err = client.AddTracksToPlaylist(ctx, newPlaylist.ID, newTracks...)
		if err != nil {
			log.Fatal(err)
		}
		newTracks = make([]spotify.ID, 0, len(tracks.Items))
		err = client.NextPage(ctx, tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
