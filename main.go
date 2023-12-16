package main

import (
	"bufio"
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

	ttsinput, ttsoutput, err := StartTTS(settings.Voicevox.SpeakerID)
	if err != nil {
		fmt.Println("Failed to StartTTS.", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input text...")
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
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

		fmt.Println("Input text...")
	}
}