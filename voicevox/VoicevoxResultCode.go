package voicevox

type ResultCode int32

const (
	/**
	 * 成功
	 */
	RESULT_OK ResultCode = iota
	/**
	* open_jtalk辞書ファイルが読み込まれていない
	 */
	RESULT_NOT_LOADED_OPENJTALK_DICT_ERROR
	/**
	* modelの読み込みに失敗した
	 */
	RESULT_LOAD_MODEL_ERROR
	/**
	* サポートされているデバイス情報取得に失敗した
	 */
	RESULT_GET_SUPPORTED_DEVICES_ERROR
	/**
	* GPUモードがサポートされていない
	 */
	RESULT_GPU_SUPPORT_ERROR
	/**
	* メタ情報読み込みに失敗した
	 */
	RESULT_LOAD_METAS_ERROR
	/**
	* ステータスが初期化されていない
	 */
	RESULT_UNINITIALIZED_STATUS_ERROR
	/**
	* 無効なspeaker_idが指定された
	 */
	RESULT_INVALID_SPEAKER_ID_ERROR
	/**
	* 無効なmodel_indexが指定された
	 */
	RESULT_INVALID_MODEL_INDEX_ERROR
	/**
	* 推論に失敗した
	 */
	RESULT_INFERENCE_ERROR
	/**
	* コンテキストラベル出力に失敗した
	 */
	RESULT_EXTRACT_FULL_CONTEXT_LABEL_ERROR
	/**
	* 無効なutf8文字列が入力された
	 */
	RESULT_INVALID_UTF8_INPUT_ERROR
	/**
	* aquestalk形式のテキストの解析に失敗した
	 */
	RESULT_PARSE_KANA_ERROR
	/**
	* 無効なAudioQuery
	 */
	RESULT_INVALID_AUDIO_QUERY_ERROR
)

func (r ResultCode) Error() string {
	switch r {
	/**
	 * 成功
	 */
	case RESULT_OK:
		return "RESULT_OK"
	/**
	* open_jtalk辞書ファイルが読み込まれていない
	 */
	case RESULT_NOT_LOADED_OPENJTALK_DICT_ERROR:
		return "RESULT_NOT_LOADED_OPENJTALK_DICT_ERROR"
	/**
	* modelの読み込みに失敗した
	 */
	case RESULT_LOAD_MODEL_ERROR:
		return "RESULT_LOAD_MODEL_ERROR"
	/**
	* サポートされているデバイス情報取得に失敗した
	 */
	case RESULT_GET_SUPPORTED_DEVICES_ERROR:
		return "RESULT_GET_SUPPORTED_DEVICES_ERROR"
	/**
	* GPUモードがサポートされていない
	 */
	case RESULT_GPU_SUPPORT_ERROR:
		return "RESULT_GPU_SUPPORT_ERROR"
	/**
	* メタ情報読み込みに失敗した
	 */
	case RESULT_LOAD_METAS_ERROR:
		return "RESULT_LOAD_METAS_ERROR"
	/**
	* ステータスが初期化されていない
	 */
	case RESULT_UNINITIALIZED_STATUS_ERROR:
		return "RESULT_UNINITIALIZED_STATUS_ERROR"
	/**
	* 無効なspeaker_idが指定された
	 */
	case RESULT_INVALID_SPEAKER_ID_ERROR:
		return "RESULT_INVALID_SPEAKER_ID_ERROR"
	/**
	* 無効なmodel_indexが指定された
	 */
	case RESULT_INVALID_MODEL_INDEX_ERROR:
		return "RESULT_INVALID_MODEL_INDEX_ERROR"
	/**
	* 推論に失敗した
	 */
	case RESULT_INFERENCE_ERROR:
		return "RESULT_INFERENCE_ERROR"
	/**
	* コンテキストラベル出力に失敗した
	 */
	case RESULT_EXTRACT_FULL_CONTEXT_LABEL_ERROR:
		return "RESULT_EXTRACT_FULL_CONTEXT_LABEL_ERROR"
	/**
	* 無効なutf8文字列が入力された
	 */
	case RESULT_INVALID_UTF8_INPUT_ERROR:
		return "RESULT_INVALID_UTF8_INPUT_ERROR"
	/**
	* aquestalk形式のテキストの解析に失敗した
	 */
	case RESULT_PARSE_KANA_ERROR:
		return "RESULT_PARSE_KANA_ERROR"
	/**
	* 無効なAudioQuery
	 */
	case RESULT_INVALID_AUDIO_QUERY_ERROR:
		return "RESULT_INVALID_AUDIO_QUERY_ERROR"
	}
	return ""
}
