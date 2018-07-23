package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

// SpotifyDevice represents a configured Spotify playback device.
type SpotifyDevice struct {
	deviceID    spotify.ID
	client      *spotify.Client
	playOpts    *spotify.PlayOptions
	playerState *spotify.PlayerState
}

// NewSpotifyDevice creates a new object associated to the Spotify
// device with the given device ID.
func NewSpotifyDevice(deviceID string) *SpotifyDevice {
	ret := &SpotifyDevice{
		deviceID: spotify.ID(deviceID),
		playOpts: &spotify.PlayOptions{},
	}
	ret.playOpts.DeviceID = &ret.deviceID
	return ret
}

func (t *SpotifyDevice) refreshPlayerState() {
	if t.client == nil {
		return
	}
	state, err := t.client.PlayerState()
	if err != nil {
		log.Println(err)
		t.playerState = nil
		return
	}
	t.playerState = state
}

// IsPlaying determines if the device is currently playing a song.
func (t *SpotifyDevice) IsPlaying() bool {
	t.refreshPlayerState()
	return (t.playerState != nil) &&
		(t.playerState.Device.ID == t.deviceID) &&
		(t.playerState.CurrentlyPlaying.Playing)
}

// Play tries to play on the configured device. This may not be allowed.
func (t *SpotifyDevice) Play() {
	if t.client == nil {
		return
	}
	log.Println("Playing on Spotify device")
	if err := t.client.PlayOpt(t.playOpts); err != nil {
		log.Println(err)
	}
}

// Pause pauses the configured device.
func (t *SpotifyDevice) Pause() {
	if t.client == nil {
		return
	}
	log.Println("Pausing Spotify device")
	if err := t.client.PauseOpt(t.playOpts); err != nil {
		log.Println(err)
	}
}

// NextTrack moves to the next track.
func (t *SpotifyDevice) NextTrack() {
	if t.client == nil {
		return
	}
	log.Println("Next track on Spotify device")
	if err := t.client.NextOpt(t.playOpts); err != nil {
		log.Println(err)
	}
}

// PrevTrack moves to the previous track.
func (t *SpotifyDevice) PrevTrack() {
	if t.client == nil {
		return
	}
	log.Println("Previous track on Spotify device")
	if err := t.client.PreviousOpt(t.playOpts); err != nil {
		log.Println(err)
	}
}

// ListDevices logs the available spotify devices and their IDs.
func (t *SpotifyDevice) ListDevices() {
	if t.client == nil {
		return
	}
	devices, err := t.client.PlayerDevices()
	if err != nil {
		log.Println(err)
		return
	}
	for _, device := range devices {
		fmt.Printf("%s %s\n", device.ID, device.Name)
	}
}
