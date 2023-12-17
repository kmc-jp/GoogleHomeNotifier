package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func GetWavLength(FilePath string) (float64, error) {
	file, err := os.Open(FilePath)
	if err != nil {
		return 0, fmt.Errorf("Open: %v", err)
	}
	defer file.Close()

	var sampleRate uint32
	var dataSize uint32

	file.Seek(24, 0)

	binary.Read(file, binary.LittleEndian, &sampleRate)

	file.Seek(40, 0)

	binary.Read(file, binary.LittleEndian, &dataSize)

	return float64(dataSize) / float64(sampleRate*2), nil
}
