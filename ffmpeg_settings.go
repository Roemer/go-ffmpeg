package goffmpeg

type FFmpegSettings struct {
	Overwrite            bool
	StatsPeriod          int
	ExecutablePath       string
	FfprobePath          string
	ShowCommandInConsole bool
	LogLevel             LogLevel
	LogMessageAction     func(string)
}

func DefaultFFmpegSettings() *FFmpegSettings {
	return &FFmpegSettings{
		Overwrite:      true,
		StatsPeriod:    30,
		ExecutablePath: "ffmpeg",
		FfprobePath:    "ffprobe",
		LogLevel:       LogLevelError,
	}
}

func (s *FFmpegSettings) SetShowCommandInConsole(show bool) *FFmpegSettings {
	s.ShowCommandInConsole = show
	return s
}

type X264Settings struct {
	Preset X264Preset
	Tune   X264Tune
	CRF    int
}

func NewX264Settings() *X264Settings {
	return &X264Settings{
		Preset: X264PresetSlow,
		Tune:   X264TuneAnimation,
		CRF:    19,
	}
}

func (s *X264Settings) SetPreset(p X264Preset) *X264Settings { s.Preset = p; return s }
func (s *X264Settings) SetTune(t X264Tune) *X264Settings     { s.Tune = t; return s }
func (s *X264Settings) SetCrf(crf int) *X264Settings         { s.CRF = crf; return s }
