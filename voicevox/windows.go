//go:build windows
// +build windows

package voicevox

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var voicevoxcoredll = syscall.MustLoadDLL(`voicevox_core.dll`)

var (
	make_default_initialize_options_proc  = voicevoxcoredll.MustFindProc("voicevox_make_default_initialize_options")
	initialize_proc                       = voicevoxcoredll.MustFindProc("voicevox_initialize")
	get_version_proc                      = voicevoxcoredll.MustFindProc("voicevox_get_version")
	load_model_proc                       = voicevoxcoredll.MustFindProc("voicevox_load_model")
	is_gpu_mode_proc                      = voicevoxcoredll.MustFindProc("voicevox_is_gpu_mode")
	is_model_loaded_proc                  = voicevoxcoredll.MustFindProc("voicevox_is_model_loaded")
	finalize_proc                         = voicevoxcoredll.MustFindProc("voicevox_finalize")
	get_metas_json_proc                   = voicevoxcoredll.MustFindProc("voicevox_get_metas_json")
	get_supported_devices_json_proc       = voicevoxcoredll.MustFindProc("voicevox_get_supported_devices_json")
	predict_duration_proc                 = voicevoxcoredll.MustFindProc("voicevox_predict_duration")
	predict_duration_data_free_proc       = voicevoxcoredll.MustFindProc("voicevox_predict_duration_data_free")
	predict_intonation_proc               = voicevoxcoredll.MustFindProc("voicevox_predict_intonation")
	predict_intonation_data_free_proc     = voicevoxcoredll.MustFindProc("voicevox_predict_intonation_data_free")
	decode_proc                           = voicevoxcoredll.MustFindProc("voicevox_decode")
	decode_data_free_proc                 = voicevoxcoredll.MustFindProc("voicevox_decode_data_free")
	make_default_audio_query_options_proc = voicevoxcoredll.MustFindProc("voicevox_make_default_audio_query_options")
	audio_query_proc                      = voicevoxcoredll.MustFindProc("voicevox_audio_query")
	make_default_synthesis_options_proc   = voicevoxcoredll.MustFindProc("voicevox_make_default_synthesis_options")
	synthesis_proc                        = voicevoxcoredll.MustFindProc("voicevox_synthesis")
	make_default_tts_options_proc         = voicevoxcoredll.MustFindProc("voicevox_make_default_tts_options")
	tts_proc                              = voicevoxcoredll.MustFindProc("voicevox_tts")
	audio_query_json_free_proc            = voicevoxcoredll.MustFindProc("voicevox_audio_query_json_free")
	wav_free_proc                         = voicevoxcoredll.MustFindProc("voicevox_wav_free")
	error_result_to_message_proc          = voicevoxcoredll.MustFindProc("voicevox_error_result_to_message")
)

// Problem: Program crashes when this function was called
// func MakeDefaultInitializeOptions() VoicevoxInitializeOptions {
// 	r1, _, _ := make_default_initialize_options_proc.Call()
// 	return *(*VoicevoxInitializeOptions)(unsafe.Pointer(&r1))
// }

func Initialize(options VoicevoxInitializeOptions) error {
	var conv_options = struct {
		AccelerationMode VoicevoxAccelerationMode
		CpuNumThreads    uint16
		LoadAllModels    bool
		OpenJtalkDictDir []byte
	}{
		AccelerationMode: options.AccelerationMode,
		CpuNumThreads:    options.CpuNumThreads,
		LoadAllModels:    options.LoadAllModels,
		OpenJtalkDictDir: append([]byte(options.OpenJtalkDictDir), 0x00),
	}

	r1, _, _ := initialize_proc.Call(uintptr(unsafe.Pointer(&conv_options.AccelerationMode)))
	if r1 != 0 {
		return fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}
	return nil
}

func GetVersion() string {
	r1, _, _ := get_version_proc.Call()
	return UTF8PtrToString((*byte)(unsafe.Pointer(r1)))
}

func LoadModel(speaker_id uint32) error {
	r1, _, _ := load_model_proc.Call(uintptr(speaker_id))
	if r1 != 0 {
		return fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}
	return nil
}

func IsGPUMode() bool {
	r1, _, _ := is_gpu_mode_proc.Call()
	return r1 == 1
}

func IsModelLoaded(speaker_id uint32) bool {
	r1, _, _ := is_model_loaded_proc.Call(uintptr(speaker_id))
	return r1 == 1
}

func Finalize() {
	finalize_proc.Call()
}

func GetMetasJson() string {
	r1, _, _ := get_metas_json_proc.Call()
	return UTF8PtrToString((*byte)(unsafe.Pointer(r1)))
}

func GetSupportedDevicesJson() string {
	r1, _, _ := get_supported_devices_json_proc.Call()
	return UTF8PtrToString((*byte)(unsafe.Pointer(r1)))
}

func PredictDuration(
	phoneme_vector []int64,
	speaker_id uint32,
) (durations []float32, err error) {
	var output_predict_duration_data_length uintptr
	var output_predict_duration_data *float32
	if len(phoneme_vector) == 0 {
		return nil, fmt.Errorf("invalid PhonemeVector size")
	}
	r1, _, _ := predict_duration_proc.Call(
		uintptr(len(phoneme_vector)),
		uintptr(unsafe.Pointer(&phoneme_vector[0])),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&output_predict_duration_data_length)),
		uintptr(unsafe.Pointer(&output_predict_duration_data)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	var rawDurations = unsafe.Slice(output_predict_duration_data, output_predict_duration_data_length)

	durations = make([]float32, output_predict_duration_data_length)

	copy(durations, rawDurations)

	predictDurationDataFree(rawDurations)

	return durations, nil
}

func predictDurationDataFree(durations []float32) {
	predict_duration_data_free_proc.Call(uintptr(unsafe.Pointer(&durations[0])))
}

func PredictIntonation(
	vowel_phoneme_vector []int64,
	consonant_phoneme_vector []int64,
	start_accent_vector []int64,
	end_accent_vector []int64,
	start_accent_phrase_vector []int64,
	end_accent_phrase_vector []int64,
	speaker_id uint32,
) (intonations []float32, err error) {
	var output_predict_intonation_data_length uintptr
	var output_predict_intonation_data *float32
	r1, _, _ := predict_intonation_proc.Call(
		uintptr(len(vowel_phoneme_vector)),
		uintptr(unsafe.Pointer(&vowel_phoneme_vector[0])),
		uintptr(unsafe.Pointer(&consonant_phoneme_vector[0])),
		uintptr(unsafe.Pointer(&start_accent_vector[0])),
		uintptr(unsafe.Pointer(&end_accent_vector[0])),
		uintptr(unsafe.Pointer(&start_accent_phrase_vector[0])),
		uintptr(unsafe.Pointer(&end_accent_phrase_vector[0])),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&output_predict_intonation_data_length)),
		uintptr(unsafe.Pointer(&output_predict_intonation_data)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	var rawIntonations = unsafe.Slice(output_predict_intonation_data, output_predict_intonation_data_length)

	intonations = make([]float32, output_predict_intonation_data_length)

	copy(intonations, rawIntonations)
	predictIntonationDataFree(rawIntonations)

	return intonations, nil
}

func predictIntonationDataFree(intonations []float32) {
	predict_intonation_data_free_proc.Call(uintptr(unsafe.Pointer(&intonations[0])))
}

func Decode(
	f0 []float32,
	phoneme_vector []float32,
	speaker_id uint32,
) ([]float32, error) {
	var length = len(f0)
	var phoneme_size = len(phoneme_vector) / length

	var output_decode_data_length uintptr
	var output_decode_data *float32
	r1, _, _ := decode_proc.Call(
		uintptr(length),
		uintptr(phoneme_size),
		uintptr(unsafe.Pointer(&f0[0])),
		uintptr(unsafe.Pointer(&phoneme_vector[0])),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&output_decode_data_length)),
		uintptr(unsafe.Pointer(&output_decode_data)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	var rawDecode_data = unsafe.Slice(output_decode_data, output_decode_data_length)

	var decode_data = make([]float32, output_decode_data_length)

	copy(decode_data, rawDecode_data)
	decodeDataFree(rawDecode_data)

	return decode_data, nil
}

func decodeDataFree(decode_data []float32) {
	decode_data_free_proc.Call(uintptr(unsafe.Pointer(&decode_data[0])))
}

// TODO: implement voicevox_make_default_audio_query_options

func AudioQuery(
	text string,
	speaker_id uint32,
	options VoicevoxAudioQueryOptions,
) (string, error) {
	var rawoutput string
	var conv_text = append([]byte(text), 0x00)
	r1, _, _ := audio_query_proc.Call(
		uintptr(unsafe.Pointer(&conv_text[0])),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&options.Kana)),
		uintptr(unsafe.Pointer(&rawoutput)),
	)
	if r1 != 0 {
		return "", fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	var output = strings.Clone(rawoutput)
	audioQueryJsonFree(rawoutput)

	return output, nil
}

// TODO: implement voicevox_make_default_synthesis_options

func Synthesis(
	audio_query string,
	speaker_id uint32,
	options VoicevoxSynthesisOptions,
) ([]byte, error) {
	var output_wav *byte
	var output_wav_length uintptr
	r1, _, _ := synthesis_proc.Call(
		uintptr(unsafe.Pointer(&audio_query)),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&options.EnableInterrogativeUpspeak)),
		uintptr(unsafe.Pointer(&output_wav_length)),
		uintptr(unsafe.Pointer(&output_wav)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	output := unsafe.Slice(output_wav, output_wav_length)
	return output, nil
}

func TTS(
	text string,
	speaker_id uint32,
	options VoicevoxTtsOptions,
) ([]byte, error) {
	var output_wav *byte
	var output_wav_length uintptr
	var conv_text = append([]byte(text), 0x00)
	r1, _, _ := tts_proc.Call(
		uintptr(unsafe.Pointer(&conv_text[0])),
		uintptr(speaker_id),
		uintptr(unsafe.Pointer(&options.Kana)),
		uintptr(unsafe.Pointer(&output_wav_length)),
		uintptr(unsafe.Pointer(&output_wav)),
	)
	if r1 != 0 {
		return nil, fmt.Errorf(ErrorResultToMessage(ResultCode(r1)))
	}

	output := unsafe.Slice(output_wav, output_wav_length)
	return output, nil
}

func audioQueryJsonFree(audioQueryJson string) {
	audio_query_json_free_proc.Call(uintptr(unsafe.Pointer(&audioQueryJson)))
}

func WavFree(output_wav []byte) {
	wav_free_proc.Call(uintptr(unsafe.Pointer(&output_wav[0])))
}

func ErrorResultToMessage(result ResultCode) string {
	r1, _, _ := error_result_to_message_proc.Call(uintptr(result))
	return UTF8PtrToString((*byte)(unsafe.Pointer(r1)))
}

func UTF8PtrToString(p *byte) string {
	if p == nil {
		return ""
	}

	var char byte
	var chars = []byte{}

	for i := 0; ; i++ {
		char = *(*byte)(unsafe.Pointer(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(byte(0))*uintptr(i))))
		// null char
		if char == 0 {
			break
		}
		chars = append(chars, char)
	}

	return string(chars)
}
