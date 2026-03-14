package goffmpeg

type LogLevel string

const (
	LogLevelQuiet   LogLevel = "quiet"
	LogLevelPanic   LogLevel = "panic"
	LogLevelFatal   LogLevel = "fatal"
	LogLevelError   LogLevel = "error"
	LogLevelWarning LogLevel = "warning"
	LogLevelInfo    LogLevel = "info"
	LogLevelVerbose LogLevel = "verbose"
	LogLevelDebug   LogLevel = "debug"
)

type StreamType string

const (
	StreamTypeVideo       StreamType = "v"
	StreamTypeAudio       StreamType = "a"
	StreamTypeSubtitles   StreamType = "s"
	StreamTypeData        StreamType = "d"
	StreamTypeAttachments StreamType = "t"
)

type ScalerFlags string

const (
	ScalerFlagsBicubic  ScalerFlags = "bicubic"
	ScalerFlagsBilinear ScalerFlags = "bilinear"
	ScalerFlagsLanczos  ScalerFlags = "lanczos"
)

type ColorFormat string

type X264Preset string

const (
	X264PresetUltrafast X264Preset = "ultrafast"
	X264PresetSuperfast X264Preset = "superfast"
	X264PresetVeryfast  X264Preset = "veryfast"
	X264PresetFaster    X264Preset = "faster"
	X264PresetFast      X264Preset = "fast"
	X264PresetMedium    X264Preset = "medium"
	X264PresetSlow      X264Preset = "slow"
	X264PresetSlower    X264Preset = "slower"
	X264PresetVeryslow  X264Preset = "veryslow"
)

type X264Tune string

const (
	X264TuneFilm        X264Tune = "film"
	X264TuneAnimation   X264Tune = "animation"
	X264TuneGrain       X264Tune = "grain"
	X264TuneStillimage  X264Tune = "stillimage"
	X264TunePSNR        X264Tune = "psnr"
	X264TuneSSIM        X264Tune = "ssim"
	X264TuneFastDecode  X264Tune = "fastdecode"
	X264TuneZerolatency X264Tune = "zerolatency"
)

type MkvDefaultMode string

const (
	MkvDefaultModeNone        MkvDefaultMode = "none"
	MkvDefaultModeInfer       MkvDefaultMode = "infer"
	MkvDefaultModeInferNoSubs MkvDefaultMode = "infer_no_subs"
)
