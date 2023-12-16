package main

import (
	"fmt"
	"os"
	"strings"
)

var DEBUG = os.Getenv("GOOGLE_HOME_DEBUG") == "on"

func main() {
	settings, err := ReadSettings()
	if err != nil {
		fmt.Println("Failed to read settings.", err)
		return
	}

	ttsinput, ttsoutput, err := StartTTS(settings.Voicevox.SpeakerID, settings.Voicevox)
	if err != nil {
		fmt.Println("Failed to StartTTS.", err)
		return
	}

	slacktextchan, slackdonechan := StartSlack(settings.Slack)

	for text := range slacktextchan {
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		ttsinput <- text
		sound := <-ttsoutput

		if sound.Error != nil {
			fmt.Println("Failed to synthesize sound: ", sound.Error)
			continue
		}

		err = Play(sound.FilePath, settings.GoogleHome)
		if err != nil {
			fmt.Println("Failed to play sound: ", err)
			continue
		}

		os.Remove(sound.FilePath)

		slackdonechan <- true
	}
}
