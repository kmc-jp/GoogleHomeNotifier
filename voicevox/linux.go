//go:build linux
// +build linux

package voicevox

/*
#cgo LDFLAGS: -L../ -lvoicevox_core
#cgo CFLAGS: -I../
#include <stdlib.h>
#include "voicevox_core.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Initialize(options VoicevoxInitializeOptions) error {
	cOptions := C.struct_VoicevoxInitializeOptions{
		acceleration_mode:   C.VoicevoxAccelerationMode(options.AccelerationMode),
		cpu_num_threads:     C.uint16_t(options.CpuNumThreads),
		load_all_models:     C.bool(options.LoadAllModels),
		open_jtalk_dict_dir: C.CString(options.OpenJtalkDictDir),
	}
	defer C.free(unsafe.Pointer(cOptions.open_jtalk_dict_dir))

	r1 := ResultCode(C.voicevox_initialize(cOptions))
	if r1 != 0 {
		return fmt.Errorf(ErrorResultToMessage(r1))
	}

	return nil
}

func GetVersion() string {
	return C.GoString(C.voicevox_get_version())
}

func LoadModel(speakerID uint32) error {
	r1 := ResultCode(C.voicevox_load_model(C.uint32_t(speakerID)))
	if r1 != 0 {
		return fmt.Errorf(ErrorResultToMessage(r1))
	}
	return nil
}

func IsGPUMode() bool {
	return bool(C.voicevox_is_gpu_mode())
}

func IsModelLoaded(speakerID uint32) bool {
	return bool(C.voicevox_is_model_loaded(C.uint32_t(speakerID)))
}

func Finalize() {
	C.voicevox_finalize()
}

func GetMetasJSON() string {
	return C.GoString(C.voicevox_get_metas_json())
}

func PredictDuration(
	phonemeVector []int64,
	speakerID uint32,
) ([]float32, error) {
	length := uintptr(len(phonemeVector))
	cPhonemeVector := (*C.int64_t)(unsafe.Pointer(&phonemeVector[0]))

	var cOutputPredictDurationDataLength C.uintptr_t
	var cOutputPredictDurationData *C.float

	resultCode := C.voicevox_predict_duration(
		C.uintptr_t(length),
		cPhonemeVector,
		C.uint32_t(speakerID),
		&cOutputPredictDurationDataLength,
		&cOutputPredictDurationData,
	)

	if resultCode != C.int(RESULT_OK) {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(resultCode)))
	}

	lengthGo := uint(cOutputPredictDurationDataLength)

	outputPredictDurationData := unsafe.Slice((*float32)(unsafe.Pointer(cOutputPredictDurationData)), lengthGo)

	durations := make([]float32, lengthGo)
	copy(durations, outputPredictDurationData)

	predictDurationDataFree(&outputPredictDurationData[0])

	return durations, nil
}

func predictDurationDataFree(predictDurationData *float32) {
	C.voicevox_predict_duration_data_free((*C.float)(predictDurationData))
}

func PredictIntonation(
	vowelPhonemeVector,
	consonantPhonemeVector,
	startAccentVector,
	endAccentVector,
	startAccentPhraseVector,
	endAccentPhraseVector []int64,
	speakerID uint32,
) ([]float32, error) {
	length := uintptr(len(vowelPhonemeVector)) // すべてのベクトルは同じ長さである必要があります

	cVowelPhonemeVector := (*C.int64_t)(unsafe.Pointer(&vowelPhonemeVector[0]))
	cConsonantPhonemeVector := (*C.int64_t)(unsafe.Pointer(&consonantPhonemeVector[0]))
	cStartAccentVector := (*C.int64_t)(unsafe.Pointer(&startAccentVector[0]))
	cEndAccentVector := (*C.int64_t)(unsafe.Pointer(&endAccentVector[0]))
	cStartAccentPhraseVector := (*C.int64_t)(unsafe.Pointer(&startAccentPhraseVector[0]))
	cEndAccentPhraseVector := (*C.int64_t)(unsafe.Pointer(&endAccentPhraseVector[0]))

	var cOutputPredictIntonationDataLength C.uintptr_t
	var cOutputPredictIntonationData *C.float

	resultCode := C.voicevox_predict_intonation(
		C.uintptr_t(length),
		cVowelPhonemeVector,
		cConsonantPhonemeVector,
		cStartAccentVector,
		cEndAccentVector,
		cStartAccentPhraseVector,
		cEndAccentPhraseVector,
		C.uint32_t(speakerID),
		&cOutputPredictIntonationDataLength,
		&cOutputPredictIntonationData,
	)

	if resultCode != C.int(RESULT_OK) {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(resultCode)))
	}

	lengthGo := uint(cOutputPredictIntonationDataLength)
	outputPredictIntonationData := unsafe.Slice((*float32)(unsafe.Pointer(cOutputPredictIntonationData)), lengthGo)

	intonations := make([]float32, lengthGo)
	copy(intonations, outputPredictIntonationData)

	predictDurationDataFree(&outputPredictIntonationData[0])

	return outputPredictIntonationData, nil
}

func predictIntonationDataFree(predictIntonationData *float32) {
	C.voicevox_predict_intonation_data_free((*C.float)(predictIntonationData))
}

func Decode(f0, phonemeVector []float32, phonemeSize uintptr, speakerID uint32) ([]float32, error) {
	length := uintptr(len(f0)) // F0とphonemeVectorは同じ長さである必要があります
	cF0 := (*C.float)(unsafe.Pointer(&f0[0]))
	cPhonemeVector := (*C.float)(unsafe.Pointer(&phonemeVector[0]))

	var cOutputDecodeDataLength C.uintptr_t
	var cOutputDecodeData *C.float

	resultCode := C.voicevox_decode(
		C.uintptr_t(length),
		C.uintptr_t(phonemeSize),
		cF0,
		cPhonemeVector,
		C.uint32_t(speakerID),
		&cOutputDecodeDataLength,
		&cOutputDecodeData,
	)

	if resultCode != C.int(RESULT_OK) {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(resultCode)))
	}

	lengthGo := uint(cOutputDecodeDataLength)
	rawDecode_data := (*float32)(unsafe.Pointer(cOutputDecodeData))

	var decode_data = make([]float32, lengthGo)
	copy(decode_data, unsafe.Slice(rawDecode_data, lengthGo))
	decodeDataFree(rawDecode_data)

	return decode_data, nil
}

func decodeDataFree(decodeData *float32) {
	C.voicevox_decode_data_free((*C.float)(decodeData))
}

func AudioQuery(text string, speakerID uint32, options VoicevoxAudioQueryOptions) (string, error) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cOptions := C.struct_VoicevoxAudioQueryOptions{
		kana: C.bool(options.Kana),
	}

	var cOutputAudioQueryJSON *C.char

	resultCode := C.voicevox_audio_query(
		cText,
		C.uint32_t(speakerID),
		cOptions,
		&cOutputAudioQueryJSON,
	)

	if resultCode != C.int(RESULT_OK) {
		return "", fmt.Errorf(ErrorResultToMessage(ResultCode(resultCode)))
	}

	outputAudioQueryJSON := C.GoString(cOutputAudioQueryJSON)
	C.voicevox_audio_query_json_free(cOutputAudioQueryJSON)

	return outputAudioQueryJSON, nil
}

func Synthesis(audioQueryJSON string, speakerID uint32, options VoicevoxSynthesisOptions) ([]byte, error) {
	cAudioQueryJSON := C.CString(audioQueryJSON)
	defer C.free(unsafe.Pointer(cAudioQueryJSON))

	cOptions := C.struct_VoicevoxSynthesisOptions{
		enable_interrogative_upspeak: C.bool(options.EnableInterrogativeUpspeak),
	}

	var cOutputWavLength C.uintptr_t
	var cOutputWav *C.uint8_t

	resultCode := C.voicevox_synthesis(
		cAudioQueryJSON,
		C.uint32_t(speakerID),
		cOptions,
		&cOutputWavLength,
		&cOutputWav,
	)

	if resultCode != C.int(RESULT_OK) {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(resultCode)))
	}

	lengthGo := uint(cOutputWavLength)

	rawOutput := (*byte)(unsafe.Pointer(cOutputWav))
	outputWav := unsafe.Slice(rawOutput, lengthGo)

	C.voicevox_wav_free(cOutputWav)

	return outputWav, nil
}

func TTS(text string, speakerID uint32, options VoicevoxTtsOptions) ([]byte, error) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cOptions := C.struct_VoicevoxTtsOptions{
		kana:                         C.bool(options.Kana),
		enable_interrogative_upspeak: C.bool(options.EnableInterrogativeUpspeak),
	}

	var cOutputWavLength C.uintptr_t
	var cOutputWav *C.uint8_t

	resultCode := C.voicevox_tts(
		cText,
		C.uint32_t(speakerID),
		cOptions,
		&cOutputWavLength,
		&cOutputWav,
	)

	if resultCode != C.int(RESULT_OK) {
		return nil, fmt.Errorf("voicevox_tts failed: %d", resultCode)
	}

	lengthGo := uint(cOutputWavLength)

	rawOutput := (*byte)(unsafe.Pointer(cOutputWav))
	outputWav := unsafe.Slice(rawOutput, lengthGo)

	return outputWav, nil
}

func WavFree(output_wav []byte) {
	C.voicevox_wav_free((*C.uint8_t)(&output_wav[0]))
}

func ErrorResultToMessage(resultCode ResultCode) string {
	return C.GoString(C.voicevox_error_result_to_message(C.VoicevoxResultCode(resultCode)))
}
