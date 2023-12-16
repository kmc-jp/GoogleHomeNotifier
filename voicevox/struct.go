package voicevox

type VoicevoxAccelerationMode int32

const (
	VOICEVOX_ACCELERATION_MODE_AUTO VoicevoxAccelerationMode = iota
	VOICEVOX_ACCELERATION_MODE_CPU
	VOICEVOX_ACCELERATION_MODE_GPU
)

type VoicevoxInitializeOptions struct {
	/**
	 * ハードウェアアクセラレーションモード
	 */
	AccelerationMode VoicevoxAccelerationMode
	/**
	* CPU利用数を指定
	* 0を指定すると環境に合わせたCPUが利用される
	 */
	CpuNumThreads uint16
	/**
	* 全てのモデルを読み込む
	 */
	LoadAllModels bool
	/**
	* open_jtalkの辞書ディレクトリ
	 */
	OpenJtalkDictDir string
}

type VoicevoxAudioQueryOptions struct {
	/**
	 * aquestalk形式のkanaとしてテキストを解釈する
	 */
	Kana bool
}

type VoicevoxSynthesisOptions struct {
	/**
	 * 疑問文の調整を有効にする
	 */
	EnableInterrogativeUpspeak bool
}
type VoicevoxTtsOptions struct {
	/**
	 * aquestalk形式のkanaとしてテキストを解釈する
	 */
	Kana bool
}

type VoicevoxUserDictWord struct{}
