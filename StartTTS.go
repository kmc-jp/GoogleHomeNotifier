package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/TKMAX777/GoogleHomeNotifier/voicevox"
)

type TtsOutputAttr struct {
	FilePath string
	Error    error
}

func StartTTS(speakerID uint32, settings VoicevoxSetting) (chan string, chan TtsOutputAttr, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := voicevox.Initialize(voicevox.VoicevoxInitializeOptions{
		AccelerationMode: voicevox.VOICEVOX_ACCELERATION_MODE_AUTO,
		CpuNumThreads:    0,
		LoadAllModels:    false,
		OpenJtalkDictDir: "open_jtalk_dic_utf_8-1.11",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Initialize: %v", err)
	}

	err = voicevox.LoadModel(speakerID)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadModel: %v", err)
	}

	var inputchan = make(chan string, 0)

	var outputchan = make(chan TtsOutputAttr, 0)

	go func() {
		for {
			text := <-inputchan
			output, err := voicevox.TTS(text, 3, voicevox.VoicevoxTtsOptions{Kana: false})
			if err != nil {
				outputchan <- TtsOutputAttr{Error: fmt.Errorf("TTS: %v", err)}
				continue
			}

			f, err := os.CreateTemp("", "GoogleHomeSound*.wav")
			if err != nil {
				outputchan <- TtsOutputAttr{Error: fmt.Errorf("Create: %v", err)}
				continue
			}

			f.Write(output)
			f.Close()

			voicevox.WavFree(output)

			outputchan <- TtsOutputAttr{FilePath: f.Name()}
		}
	}()

	return inputchan, outputchan, nil
}
