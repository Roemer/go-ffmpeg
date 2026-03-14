package goffmpeg

import "fmt"

type encodingBase struct {
	StreamIndex int
}

func newEncodingBase() encodingBase {
	return encodingBase{StreamIndex: -1}
}

func (e *encodingBase) buildParam(encodingType, codecName string, extra map[string]string) []string {
	params := []string{}
	codec := fmt.Sprintf("-c:%s", encodingType)
	if e.StreamIndex >= 0 {
		codec += fmt.Sprintf(":%d", e.StreamIndex)
	}
	params = append(params, fmt.Sprintf("%s %s", codec, codecName))
	for k, v := range extra {
		key := k
		if e.StreamIndex >= 0 {
			key += fmt.Sprintf(":%d", e.StreamIndex)
		}
		params = append(params, fmt.Sprintf("%s %s", key, v))
	}
	return params
}

//////////
// VideoCopyEncoding
//////////

type VideoCopyEncoding struct{ encodingBase }

func NewVideoCopyEncoding() *VideoCopyEncoding {
	return &VideoCopyEncoding{encodingBase: newEncodingBase()}
}

func (e *VideoCopyEncoding) SetIndex(i int) *VideoCopyEncoding { e.StreamIndex = i; return e }

func (e *VideoCopyEncoding) GetParameterStrings() []string {
	return e.buildParam("v", "copy", nil)
}

//////////
// FFV1VideoEncoding
//////////

type FFV1VideoEncoding struct{ encodingBase }

func NewFFV1VideoEncoding() *FFV1VideoEncoding {
	return &FFV1VideoEncoding{encodingBase: newEncodingBase()}
}

func (e *FFV1VideoEncoding) SetIndex(i int) *FFV1VideoEncoding { e.StreamIndex = i; return e }

func (e *FFV1VideoEncoding) GetParameterStrings() []string {
	return e.buildParam("v", "ffv1", map[string]string{
		"-level":   "1",
		"-coder":   "1",
		"-context": "1",
		"-g":       "1",
	})
}

//////////
// X264VideoEncoding
//////////

type X264VideoEncoding struct {
	encodingBase
	Preset      X264Preset
	Tune        X264Tune
	Crf         float64
	Bitrate     *int
	Pass        int
	PassLogFile string
}

func NewX264VideoEncoding() *X264VideoEncoding {
	return &X264VideoEncoding{
		encodingBase: newEncodingBase(),
		Preset:       X264PresetSlow,
		Tune:         X264TuneAnimation,
		Crf:          19,
	}
}

func (e *X264VideoEncoding) SetIndex(i int) *X264VideoEncoding          { e.StreamIndex = i; return e }
func (e *X264VideoEncoding) SetPreset(p X264Preset) *X264VideoEncoding  { e.Preset = p; return e }
func (e *X264VideoEncoding) SetTune(t X264Tune) *X264VideoEncoding      { e.Tune = t; return e }
func (e *X264VideoEncoding) SetCrf(crf float64) *X264VideoEncoding      { e.Crf = crf; return e }
func (e *X264VideoEncoding) SetBitrate(b *int) *X264VideoEncoding       { e.Bitrate = b; return e }
func (e *X264VideoEncoding) FirstPass() *X264VideoEncoding              { e.Pass = 1; return e }
func (e *X264VideoEncoding) SecondPass() *X264VideoEncoding             { e.Pass = 2; return e }
func (e *X264VideoEncoding) SetPassLogFile(f string) *X264VideoEncoding { e.PassLogFile = f; return e }

func (e *X264VideoEncoding) FromSettings(s *X264Settings) *X264VideoEncoding {
	e.Preset = s.Preset
	e.Tune = s.Tune
	e.Crf = float64(s.CRF)
	return e
}

func (e *X264VideoEncoding) GetParameterStrings() []string {
	extra := map[string]string{
		"-preset":    string(e.Preset),
		"-tune":      string(e.Tune),
		"-profile:v": "high",
		"-level":     "4.1",
		"-pix_fmt":   "yuv420p",
	}
	if e.Bitrate != nil {
		extra["-b:v"] = fmt.Sprintf("%dk", *e.Bitrate)
		if e.Pass > 0 {
			extra["-pass"] = fmt.Sprintf("%d", e.Pass)
			if e.PassLogFile != "" {
				extra["-passlogfile"] = fmt.Sprintf("%q", e.PassLogFile)
			}
		}
	} else {
		extra["-crf"] = fmt.Sprintf("%.2f", e.Crf)
	}
	return e.buildParam("v", "libx264", extra)
}

//////////
// AudioCopyEncoding
//////////

type AudioCopyEncoding struct{ encodingBase }

func NewAudioCopyEncoding() *AudioCopyEncoding {
	return &AudioCopyEncoding{encodingBase: newEncodingBase()}
}

func (e *AudioCopyEncoding) SetIndex(i int) *AudioCopyEncoding { e.StreamIndex = i; return e }

func (e *AudioCopyEncoding) GetParameterStrings() []string {
	return e.buildParam("a", "copy", nil)
}

//////////
// AacAudioEncoding
//////////

type AacAudioEncoding struct {
	encodingBase
	Bitrate       int
	ChannelLayout string
}

func NewAacAudioEncoding() *AacAudioEncoding {
	return &AacAudioEncoding{encodingBase: newEncodingBase(), Bitrate: 160}
}

func (e *AacAudioEncoding) SetIndex(i int) *AacAudioEncoding   { e.StreamIndex = i; return e }
func (e *AacAudioEncoding) SetBitrate(b int) *AacAudioEncoding { e.Bitrate = b; return e }
func (e *AacAudioEncoding) SetChannelLayout(l string) *AacAudioEncoding {
	e.ChannelLayout = l
	return e
}

func (e *AacAudioEncoding) GetParameterStrings() []string {
	extra := map[string]string{"-b:a": fmt.Sprintf("%dk", e.Bitrate)}
	if e.ChannelLayout != "" {
		extra["-channel_layout"] = fmt.Sprintf("%q", e.ChannelLayout)
	}
	return e.buildParam("a", "aac", extra)
}

//////////
// Ac3AudioEncoding
//////////

type Ac3AudioEncoding struct {
	encodingBase
	Bitrate       int
	ChannelLayout string
}

func NewAc3AudioEncoding() *Ac3AudioEncoding {
	return &Ac3AudioEncoding{encodingBase: newEncodingBase(), Bitrate: 192}
}

func (e *Ac3AudioEncoding) SetIndex(i int) *Ac3AudioEncoding   { e.StreamIndex = i; return e }
func (e *Ac3AudioEncoding) SetBitrate(b int) *Ac3AudioEncoding { e.Bitrate = b; return e }
func (e *Ac3AudioEncoding) SetChannelLayout(l string) *Ac3AudioEncoding {
	e.ChannelLayout = l
	return e
}

func (e *Ac3AudioEncoding) GetParameterStrings() []string {
	extra := map[string]string{"-b:a": fmt.Sprintf("%dk", e.Bitrate)}
	if e.ChannelLayout != "" {
		extra["-channel_layout"] = fmt.Sprintf("%q", e.ChannelLayout)
	}
	return e.buildParam("a", "ac3", extra)
}

//////////
// Mp3AudioEncoding
//////////

type Mp3AudioEncoding struct {
	encodingBase
	Bitrate int
}

func NewMp3AudioEncoding() *Mp3AudioEncoding {
	return &Mp3AudioEncoding{encodingBase: newEncodingBase(), Bitrate: 192}
}

func (e *Mp3AudioEncoding) SetIndex(i int) *Mp3AudioEncoding   { e.StreamIndex = i; return e }
func (e *Mp3AudioEncoding) SetBitrate(b int) *Mp3AudioEncoding { e.Bitrate = b; return e }

func (e *Mp3AudioEncoding) GetParameterStrings() []string {
	return e.buildParam("a", "libmp3lame", map[string]string{"-b:a": fmt.Sprintf("%dk", e.Bitrate)})
}
