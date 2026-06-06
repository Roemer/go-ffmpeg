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

type ErrorDetectionFlag string

const (
	// verify embedded CRCs
	ErrorDetectionFlagCRCCheck ErrorDetectionFlag = "crccheck"
	// detect bitstream specification deviations
	ErrorDetectionFlagBitstream ErrorDetectionFlag = "bitstream"
	// detect improper bitstream length
	ErrorDetectionFlagBuffer ErrorDetectionFlag = "buffer"
	// abort decoding on minor error detection
	ErrorDetectionFlagExplode ErrorDetectionFlag = "explode"
	// ignore decoding errors, and continue decoding.
	// This is useful if you want to analyze the content of a video and thus want everything to be decoded no matter what.
	// This option will not result in a video that is pleasing to watch in case of errors.
	ErrorDetectionFlagIgnoreErr ErrorDetectionFlag = "ignore_err"
	// consider things that violate the spec and have not been seen in the wild as errors
	ErrorDetectionFlagCareful ErrorDetectionFlag = "careful"
	// consider all spec non compliancies as errors
	ErrorDetectionFlagCompliant ErrorDetectionFlag = "compliant"
	// consider things that a sane encoder should not do as an error
	ErrorDetectionFlagAggressive ErrorDetectionFlag = "aggressive"
)

type StreamType string

const (
	StreamTypeNone        StreamType = ""
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
