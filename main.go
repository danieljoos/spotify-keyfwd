package main

import (
	"flag"
)

func main() {
	flagDevices := flag.Bool("devices", false, "List available Spotify devices")
	flag.Parse()

	configFile := flag.Arg(0)
	if configFile == "" {
		configFile = "spotify-keyfwd.json"
	}

	config := LoadConfiguration(configFile)
	player := NewPlayer(config)
	authenticateSpotifyDevice(player.device, config)

	if *flagDevices {
		player.device.ListDevices()
	} else {
		player.Start()
		defer player.Stop()
		<-player.terminate
	}
}
