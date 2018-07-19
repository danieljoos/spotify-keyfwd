package main

// Player is acting as controller for this application.
// It receives keyboard capture events and forwards them to its
// spotify device.
type Player struct {
	keyboardCapture *KeyboardCapture
	device          *SpotifyDevice
	terminate       chan struct{}
}

type playerAction int

const (
	playerActionNone playerAction = iota
	playerActionPlay
	playerActionPause
	playerActionNext
	playerActionPrev
)

// NewPlayer creates a new player object using the given configuration.
func NewPlayer(config *Configuration) *Player {
	return &Player{
		keyboardCapture: NewKeyboardCapture(playerCaptureKeys),
		device:          NewSpotifyDevice(config.DeviceID),
		terminate:       make(chan struct{}),
	}
}

// Start causes the player object to start listening for keyboard events.
func (t *Player) Start() {
	go t.keyboardCapture.SyncReceive()
	go func() {
		for {
			select {
			case key := <-t.keyboardCapture.KeyPressed:
				t.onKeyPress(key)
			case <-t.terminate:
				return
			}
		}
	}()
}

// Stop stops listening for keyboard events and sets the termination
// signal of the player object.
func (t *Player) Stop() {
	t.keyboardCapture.Stop()
	t.terminate <- struct{}{}
}

// IsPlaying returns true if the player's spotify device is currently
// playing a song. Otherwise false.
func (t *Player) IsPlaying() bool {
	return t.device.IsPlaying()
}

func (t *Player) onKeyPress(key int) {
	switch t.keyToAction(key) {
	case playerActionPlay:
		t.device.Play()
	case playerActionPause:
		t.device.Pause()
	case playerActionNext:
		t.device.NextTrack()
	case playerActionPrev:
		t.device.PrevTrack()
	}
}
