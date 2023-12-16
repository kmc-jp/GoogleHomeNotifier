package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type Setting struct {
	Voicevox   VoicevoxSetting   `yaml:"Voicevox"`
	GoogleHome GoogleHomeSetting `yaml:"GoogleHome"`
}

type VoicevoxSetting struct {
	SpeakerID uint32 `yaml:"SpeakerID"`
}

type GoogleHomeSetting struct {
	DeviceName  string `yaml:"DeviceName"`
	Device      string `yaml:"Device"`
	Iface       string `yaml:"Iface"`
	ForceDetach bool   `yaml:"Detach"`
	Detach      bool   `yaml:"Detach"`
	Addr        string `yaml:"Addr"`
	Port        int    `yaml:"Port"`
	UUID        string `yaml:"UUID"`
}

func ReadSettings() (*Setting, error) {
	var yamlRootPath = "settings"

	dir, err := os.ReadDir(yamlRootPath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDir")
	}

	var yamlBinary = []byte{}
	for _, f := range dir {
		if f.IsDir() || !(strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
			continue
		}

		var yamlFilePath = filepath.Join(yamlRootPath, f.Name())
		b, err := os.ReadFile(yamlFilePath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("ReadFile: %s", yamlFilePath))
		}

		// check format
		var us Setting
		err = yaml.Unmarshal(b, &us)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("UnmarshalSettings: %s\n%s", yamlFilePath, err.Error()))
		}

		yamlBinary = append(yamlBinary, b...)
		yamlBinary = append(yamlBinary, '\n')
	}

	var us Setting
	err = yaml.Unmarshal(yamlBinary, &us)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return &us, nil
}
