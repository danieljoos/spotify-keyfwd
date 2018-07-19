package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/zmb3/spotify"
)

const (
	authCallbackAddressFmt = "localhost:%d"
	authCallbackURLFmt     = "http://localhost:%d/callback"
	authResponse           = "<html><body>Browser window can be closed...<script>window.close()</script></body></html>"
)

func authenticateSpotifyDevice(device *SpotifyDevice, config *Configuration) {
	authCallbackAddress := fmt.Sprintf(authCallbackAddressFmt, config.AuthHTTPPort)
	authCallbackURL := fmt.Sprintf(authCallbackURLFmt, config.AuthHTTPPort)
	authenticator := spotify.NewAuthenticator(authCallbackURL,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState)
	authenticator.SetAuthInfo(config.SpotifyClientID, config.SpotifySecretKey)
	authenticated := make(chan struct{})
	authState := generateAuthState()

	// Prepare HTTP auth endpoint.
	// After the authentication, Spotify will redirect the browser
	// to this endpoint.
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    authCallbackAddress,
		Handler: mux,
	}
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := authenticator.Token(authState, r)
		if err != nil {
			http.Error(w, "Failed to get token", http.StatusForbidden)
			log.Println(err)
			return
		}
		if state := r.FormValue("state"); state != authState {
			http.NotFound(w, r)
			log.Println("State mismatch")
			return
		}
		client := authenticator.NewClient(token)
		device.client = &client
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, authResponse)
		authenticated <- struct{}{}
	})

	// Start the authentication process
	log.Println("Performing Spotify API authentication")
	go server.ListenAndServe()
	url := authenticator.AuthURL(authState)
	browser.OpenURL(url)
	<-authenticated
	server.Shutdown(context.Background())
	log.Println("Successfully authenticated!")
}

func generateAuthState() string {
	const glyphs = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, 10)
	for i := range res {
		res[i] = glyphs[rand.Intn(len(glyphs))]
	}
	return string(res)
}
