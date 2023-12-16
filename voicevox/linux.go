//go:build linux
// +build linux

package voicevox

/*
#cgo LDFLAGS: -L. -lvoicevox_core
#include "voicevox_core.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// ここにラッパー関数を定義

func Initialize(options VoicevoxInitializeOptions) ResultCode {
	cOptions := C.struct_VoicevoxInitializeOptions{
		acceleration_mode:   C.VoicevoxAccelerationMode(options.AccelerationMode),
		cpu_num_threads:     C.uint16_t(options.CpuNumThreads),
		load_all_models:     C.bool(options.LoadAllModels),
		open_jtalk_dict_dir: C.CString(options.OpenJtalkDictDir),
	}
	defer C.free(unsafe.Pointer(cOptions.open_jtalk_dict_dir))

	return ResultCode(C.voicevox_initialize(cOptions))
}

func GetVersion() string {
	return C.GoString(C.voicevox_get_version())
}

func LoadModel(speakerID uint32) ResultCode {
	return ResultCode(C.voicevox_load_model(C.uint32_t(speakerID)))
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

func PredictDuration(phonemeVector []int64, speakerID uint32) (ResultCode, []float32, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), nil, fmt.Errorf("voicevox_predict_duration failed: %d", resultCode)
	}

	lengthGo := uint(cOutputPredictDurationDataLength)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOutputPredictDurationData)),
		Len:  int(lengthGo),
		Cap:  int(lengthGo),
	}
	outputPredictDurationData := *(*[]float32)(unsafe.Pointer(&sliceHeader))

	return ResultCode(resultCode), outputPredictDurationData, nil
}

func PredictDurationDataFree(predictDurationData *float32) {
	C.voicevox_predict_duration_data_free((*C.float)(predictDurationData))
}

func PredictIntonation(vowelPhonemeVector, consonantPhonemeVector, startAccentVector, endAccentVector, startAccentPhraseVector, endAccentPhraseVector []int64, speakerID uint32) (ResultCode, []float32, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), nil, fmt.Errorf("voicevox_predict_intonation failed: %d", resultCode)
	}

	lengthGo := uint(cOutputPredictIntonationDataLength)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOutputPredictIntonationData)),
		Len:  int(lengthGo),
		Cap:  int(lengthGo),
	}
	outputPredictIntonationData := *(*[]float32)(unsafe.Pointer(&sliceHeader))

	return ResultCode(resultCode), outputPredictIntonationData, nil
}

func PredictIntonationDataFree(predictIntonationData *float32) {
	C.voicevox_predict_intonation_data_free((*C.float)(predictIntonationData))
}

func Decode(f0, phonemeVector []float32, phonemeSize uintptr, speakerID uint32) (ResultCode, []float32, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), nil, fmt.Errorf("voicevox_decode failed: %d", resultCode)
	}

	lengthGo := uint(cOutputDecodeDataLength)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOutputDecodeData)),
		Len:  int(lengthGo),
		Cap:  int(lengthGo),
	}
	outputDecodeData := *(*[]float32)(unsafe.Pointer(&sliceHeader))

	return ResultCode(resultCode), outputDecodeData, nil
}

func DecodeDataFree(decodeData *float32) {
	C.voicevox_decode_data_free((*C.float)(decodeData))
}

func AudioQuery(text string, speakerID uint32, options VoicevoxAudioQueryOptions) (ResultCode, string, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), "", fmt.Errorf("voicevox_audio_query failed: %d", resultCode)
	}

	outputAudioQueryJSON := C.GoString(cOutputAudioQueryJSON)
	C.voicevox_audio_query_json_free(cOutputAudioQueryJSON)

	return ResultCode(resultCode), outputAudioQueryJSON, nil
}

func AudioQueryJSONFree(audioQueryJSON *char) {
	C.voicevox_audio_query_json_free((*C.char)(audioQueryJSON))
}

func Synthesis(audioQueryJSON string, speakerID uint32, options VoicevoxSynthesisOptions) (ResultCode, []byte, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), nil, fmt.Errorf("voicevox_synthesis failed: %d", resultCode)
	}

	lengthGo := uint(cOutputWavLength)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOutputWav)),
		Len:  int(lengthGo),
		Cap:  int(lengthGo),
	}
	outputWav := *(*[]byte)(unsafe.Pointer(&sliceHeader))
	C.voicevox_wav_free(cOutputWav)

	return ResultCode(resultCode), outputWav, nil
}

func TTS(text string, speakerID uint32, options VoicevoxTtsOptions) (ResultCode, []byte, error) {
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

	if resultCode != RESULT_OK {
		return ResultCode(resultCode), nil, fmt.Errorf("voicevox_tts failed: %d", resultCode)
	}

	lengthGo := uint(cOutputWavLength)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOutputWav)),
		Len:  int(lengthGo),
		Cap:  int(lengthGo),
	}
	outputWav := *(*[]byte)(unsafe.Pointer(&sliceHeader))
	C.voicevox_wav_free(cOutputWav)

	return ResultCode(resultCode), outputWav, nil
}

func WavFree(wav *uint8) {
	C.voicevox_wav_free((*C.uint8_t)(wav))
}

func ErrorResultToMessage(resultCode ResultCode) string {
	return C.GoString(C.voicevox_error_result_to_message(C.ResultCode(resultCode)))
}
