package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/godbus/dbus"
)

var (
	pauseColor = flag.String("pausecolor", "#FFFFFF", "Specify which alternate color to use when spotify is paused")
)

func main() {
	flag.Parse()

	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect to SystemBus bus:", err)
	}
	obj := conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
	switch command := os.Getenv("BLOCK_BUTTON"); command {
	case "1":
		playPause(obj)
	case "2":
		previous(obj)
	case "3":
		next(obj)
	}
	// This sleep is needed since Spotify needs time to update it's status before
	// we call the metadata-function.
	time.Sleep(time.Millisecond * 100)

	line, playStatus := metadata(obj)
	fmt.Println(line)
	if !playStatus {
		fmt.Fprint(os.Stdout, "\n", *pauseColor, "\n")
	}
}

func playPause(obj dbus.BusObject) {
	call := obj.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Error: Spotify probably not running")
		os.Exit(1)
	}
}

func next(obj dbus.BusObject) {
	call := obj.Call("org.mpris.MediaPlayer2.Player.Next", 0)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Error: Spotify probably not running")
		os.Exit(1)
	}
}

func previous(obj dbus.BusObject) {
	call := obj.Call("org.mpris.MediaPlayer2.Player.Previous", 0)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Error: Spotify probably not running")
		os.Exit(1)
	}
}

func metadata(obj dbus.BusObject) (string, bool) {
	raw, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		fmt.Println("Error: Spotify probably not running")
		os.Exit(1)
	}
	metadata := raw.Value().(map[string]dbus.Variant)
	artist := metadata["xesam:artist"].Value().([]string)[0]
	trackTitle := metadata["xesam:title"].Value().(string)

	statusRaw, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	var status bool
	if statusRaw.Value().(string) == "Paused" {
		status = false
	} else if statusRaw.Value().(string) == "Playing" {
		status = true
	}
	return artist + " - " + trackTitle, status
}
