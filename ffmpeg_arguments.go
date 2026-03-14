package goffmpeg

import (
	"fmt"
	"strings"
)

type FFmpegArguments struct {
	InputFiles         []*InputFile
	Output             IOutput
	VideoFilters       []*VideoFilter
	AudioFilters       []*AudioFilter
	VideoEncodings     []IVideoEncoding
	AudioEncodings     []IAudioEncoding
	Mappings           []*Mapping
	Dispositions       []*Disposition
	Metadatas          []*Metadata
	CopyAll            bool
	HideBanner         bool
	XError             bool
	LogLevel           LogLevel
	Stats              bool
	StatsPeriod        int
	Overwrite          bool
	MapChaptersIndex   *int
	MaxInterleaveDelta *int
	FastStart          bool
	IgnoreErrors       bool
	MapMetaData        bool
	Async              *int
	DisableVideo       bool
	DisableAudio       bool
	DisableSubtitle    bool
	DisableData        bool
	Threads            *int
	FormatFlags        string
	Aspect             string
	MkvDefaultMode     MkvDefaultMode
}

func NewFFmpegArguments() *FFmpegArguments {
	return &FFmpegArguments{
		HideBanner:     true,
		XError:         true,
		LogLevel:       LogLevelWarning,
		Stats:          true,
		StatsPeriod:    60,
		Overwrite:      true,
		MkvDefaultMode: MkvDefaultModeInferNoSubs,
	}
}

func (f *FFmpegArguments) ArgumentString() string {
	return strings.Join(f.ArgumentSlice(), " ")
}

func (f *FFmpegArguments) ArgumentSlice() []string {
	return f.buildArguments()
}

func (f *FFmpegArguments) AddInput(inputFile *InputFile) *FFmpegArguments {
	inputFile.Index = len(f.InputFiles)
	f.InputFiles = append(f.InputFiles, inputFile)
	return f
}

func (f *FFmpegArguments) AddInputPath(path string) *FFmpegArguments {
	return f.AddInput(NewInputFile(path))
}

func (f *FFmpegArguments) SetOutput(outputFile *OutputFile) *FFmpegArguments {
	f.Output = outputFile
	return f
}

func (f *FFmpegArguments) SetOutputPath(path string) *FFmpegArguments {
	return f.SetOutput(NewOutputFile(path))
}

func (f *FFmpegArguments) SetNullOutput() *FFmpegArguments {
	f.Output = &NullOutput{}
	return f
}

func (f *FFmpegArguments) AddVideoFilter(vf *VideoFilter) *FFmpegArguments {
	f.VideoFilters = append(f.VideoFilters, vf)
	return f
}

func (f *FFmpegArguments) AddAudioFilter(af *AudioFilter) *FFmpegArguments {
	f.AudioFilters = append(f.AudioFilters, af)
	return f
}

func (f *FFmpegArguments) AddVideoEncoding(ve IVideoEncoding) *FFmpegArguments {
	f.VideoEncodings = append(f.VideoEncodings, ve)
	return f
}

func (f *FFmpegArguments) AddAudioEncoding(ae IAudioEncoding) *FFmpegArguments {
	f.AudioEncodings = append(f.AudioEncodings, ae)
	return f
}

func (f *FFmpegArguments) AddMapping(m *Mapping) *FFmpegArguments {
	f.Mappings = append(f.Mappings, m)
	return f
}

func (f *FFmpegArguments) SetCopyAll(copyAll bool) *FFmpegArguments {
	f.CopyAll = copyAll
	return f
}

func (f *FFmpegArguments) SetMapChaptersIndex(idx int) *FFmpegArguments {
	f.MapChaptersIndex = &idx
	return f
}

func (f *FFmpegArguments) AddDisposition(d *Disposition) *FFmpegArguments {
	f.Dispositions = append(f.Dispositions, d)
	return f
}

func (f *FFmpegArguments) AddMetadata(m *Metadata) *FFmpegArguments {
	f.Metadatas = append(f.Metadatas, m)
	return f
}

func (f *FFmpegArguments) Disable(disableVideo, disableAudio, disableSubtitle, disableData bool) *FFmpegArguments {
	f.DisableVideo = disableVideo
	f.DisableAudio = disableAudio
	f.DisableSubtitle = disableSubtitle
	f.DisableData = disableData
	return f
}

func (f *FFmpegArguments) SetDisableVideo(v bool) *FFmpegArguments { f.DisableVideo = v; return f }
func (f *FFmpegArguments) SetDisableAudio(v bool) *FFmpegArguments { f.DisableAudio = v; return f }
func (f *FFmpegArguments) SetDisableSubtitle(v bool) *FFmpegArguments {
	f.DisableSubtitle = v
	return f
}
func (f *FFmpegArguments) SetDisableData(v bool) *FFmpegArguments { f.DisableData = v; return f }

func (f *FFmpegArguments) SetThreads(t int) *FFmpegArguments       { f.Threads = &t; return f }
func (f *FFmpegArguments) SetLogLevel(l LogLevel) *FFmpegArguments { f.LogLevel = l; return f }

func (f *FFmpegArguments) FormatFlagGenPts() *FFmpegArguments {
	f.FormatFlags = "+genpts"
	return f
}

func (f *FFmpegArguments) SetAspect(aspect string) *FFmpegArguments { f.Aspect = aspect; return f }

func (f *FFmpegArguments) SetMkvDefaultMode(m MkvDefaultMode) *FFmpegArguments {
	f.MkvDefaultMode = m
	return f
}

func (f *FFmpegArguments) SetMaxInterleaveDelta(v int) *FFmpegArguments {
	f.MaxInterleaveDelta = Ptr(v)
	return f
}

func (f *FFmpegArguments) SetFastStart(v bool) *FFmpegArguments    { f.FastStart = v; return f }
func (f *FFmpegArguments) SetIgnoreErrors(v bool) *FFmpegArguments { f.IgnoreErrors = v; return f }
func (f *FFmpegArguments) SetMapMetaData(v bool) *FFmpegArguments  { f.MapMetaData = v; return f }

func (f *FFmpegArguments) buildArguments() []string {
	var args []string

	if f.HideBanner {
		args = append(args, "-hide_banner")
	}
	if f.XError {
		args = append(args, "-xerror")
	}
	args = append(args, "-loglevel", string(f.LogLevel))
	if f.Stats {
		args = append(args, "-stats", "-stats_period", fmt.Sprintf("%d", f.StatsPeriod))
	} else {
		args = append(args, "-nostats")
	}
	if f.IgnoreErrors {
		args = append(args, "-err_detect", "ignore_err")
	}
	if f.Overwrite {
		args = append(args, "-y")
	}
	if f.DisableVideo {
		args = append(args, "-vn")
	}
	if f.DisableAudio {
		args = append(args, "-an")
	}
	if f.DisableSubtitle {
		args = append(args, "-sn")
	}
	if f.DisableData {
		args = append(args, "-dn")
	}
	if f.Threads != nil {
		args = append(args, "-threads", fmt.Sprintf("%d", *f.Threads))
		args = append(args, "-filter_threads", fmt.Sprintf("%d", *f.Threads))
		args = append(args, "-filter_complex_threads", fmt.Sprintf("%d", *f.Threads))
	}
	if f.FormatFlags != "" {
		args = append(args, "-fflags", f.FormatFlags)
	}
	for _, inp := range f.InputFiles {
		args = append(args, inp.GetParameters()...)
	}
	if f.Threads != nil {
		args = append(args, "-threads", fmt.Sprintf("%d", *f.Threads))
	}
	for _, vf := range f.VideoFilters {
		args = append(args, vf.GetParameters()...)
	}
	for _, af := range f.AudioFilters {
		args = append(args, af.GetParameters()...)
	}
	if f.CopyAll {
		args = append(args, "-c copy")
	}
	for _, m := range f.Mappings {
		args = append(args, m.GetParameters()...)
	}
	for _, ve := range f.VideoEncodings {
		args = append(args, ve.GetParameters()...)
	}
	for _, ae := range f.AudioEncodings {
		args = append(args, ae.GetParameters()...)
	}
	for _, d := range f.Dispositions {
		args = append(args, d.GetParameters()...)
	}
	for _, m := range f.Metadatas {
		args = append(args, m.GetParameters()...)
	}
	if f.MapChaptersIndex != nil {
		args = append(args, fmt.Sprintf("-map_chapters %d", *f.MapChaptersIndex))
	}
	if f.MapMetaData {
		args = append(args, "-map_metadata")
	}
	if f.MkvDefaultMode != MkvDefaultModeNone {
		args = append(args, fmt.Sprintf("-default_mode %s", f.MkvDefaultMode))
	}
	if f.MaxInterleaveDelta != nil {
		args = append(args, fmt.Sprintf("-max_interleave_delta %d", *f.MaxInterleaveDelta))
	}
	if f.FastStart {
		args = append(args, "-movflags faststart")
	}
	if f.Aspect != "" {
		args = append(args, fmt.Sprintf("-aspect %s", f.Aspect))
	}
	if f.Async != nil {
		args = append(args, fmt.Sprintf("-async %d", *f.Async))
	}
	args = append(args, "-audio_service_type ma")
	if f.Output != nil {
		args = append(args, f.Output.GetParameters()...)
	}
	return args
}
