// +build windows

package main

import (
	"github.com/TheTitanrain/w32"
)

var playerCaptureKeys = []int{
	w32.VK_PLAY,
	w32.VK_PAUSE,
	w32.VK_MEDIA_STOP,
	w32.VK_MEDIA_NEXT_TRACK,
	w32.VK_MEDIA_PREV_TRACK,
	w32.VK_MEDIA_PLAY_PAUSE,
}

func (t *Player) keyToAction(key int) playerAction {
	switch key {
	case w32.VK_PLAY:
		return playerActionPlay
	case w32.VK_PAUSE, w32.VK_MEDIA_STOP:
		return playerActionPause
	case w32.VK_MEDIA_NEXT_TRACK:
		return playerActionNext
	case w32.VK_MEDIA_PREV_TRACK:
		return playerActionPrev
	case w32.VK_MEDIA_PLAY_PAUSE:
		{
			if t.IsPlaying() {
				return playerActionPause
			}
			return playerActionPlay
		}
	}
	return playerActionNone
}
