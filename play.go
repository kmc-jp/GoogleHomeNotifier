package main

import (
	"fmt"
	"net"
	"time"

	"github.com/vishen/go-chromecast/application"
)

func Play(sound TtsOutputAttr, settings GoogleHomeSetting) error {
	applicationOptions := []application.ApplicationOption{
		application.WithDebug(DEBUG),
	}

	var iface *net.Interface
	if settings.Iface != "" {
		var err error
		if iface, err = net.InterfaceByName(settings.Iface); err != nil {
			return fmt.Errorf("unable to find interface %q: %v", settings.Iface, err)
		}
		applicationOptions = append(applicationOptions, application.WithIface(iface))
	}

	app := application.NewApplication(applicationOptions...)
	err := app.Start(settings.Addr, settings.Port)
	if err != nil {
		return fmt.Errorf("Start: %v", err)
	}

	volume := app.Volume().Level
	app.SetVolume(settings.Volume)

	err = app.Load(sound.FilePath, 0, "audio/wav", false, settings.Detach, settings.ForceDetach)
	if err != nil {
		return fmt.Errorf("Load: %v", err)
	}

	timer := time.NewTimer(time.Second * time.Duration(settings.MaxDuration))
	stopchan := make(chan bool)

	go func() {
		app.MediaWait()
		stopchan <- true
	}()

	select {
	case <-timer.C:
		app.StopMedia()
		app.SetVolume(volume)
		return fmt.Errorf("The message was too long, so it was interrupted.")
	case <-stopchan:
		app.SetVolume(volume)
		return nil
	}
}
