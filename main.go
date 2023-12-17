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

	ttsinput, ttsoutput, err := StartTTS(settings.Voicevox)
	if err != nil {
		fmt.Println("Failed to StartTTS.", err)
		return
	}

	slacktextchan, slackdonechan := StartSlack(settings.Slack)

	fmt.Println("Start waiting messages...")

	var sound TtsOutputAttr

	for text := range slacktextchan {
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		if sound.FilePath != "" {
			os.Remove(sound.FilePath)
		}

		ttsinput <- text
		sound = <-ttsoutput

		if sound.Error != nil {
			slackdonechan <- fmt.Errorf("Failed to synthesize sound: %s", sound.Error)
			continue
		}

		err = Play(sound, settings.GoogleHome)
		if err != nil {
			slackdonechan <- fmt.Errorf("Failed to play sound: %v", err)
			continue
		}

		slackdonechan <- nil
	}
}
