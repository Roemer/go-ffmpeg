package goffmpeg

import (
	"fmt"
	"strings"
)

type filterBase struct {
	filters     []string
	StreamIndex int
}

func newFilterBase() filterBase {
	return filterBase{StreamIndex: -1}
}

func (f *filterBase) addFilter(s string) {
	f.filters = append(f.filters, s)
}

func (f *filterBase) buildParameters(filterType string) []string {
	if len(f.filters) == 0 {
		return nil
	}
	idx := ""
	if f.StreamIndex >= 0 {
		idx = fmt.Sprintf(":%d", f.StreamIndex)
	}
	return []string{fmt.Sprintf("-filter:%s%s", filterType, idx), fmt.Sprintf(`"%s"`, strings.Join(f.filters, ", "))}
}

//////////
// VideoFilter
//////////

type VideoFilter struct {
	filterBase
}

func NewVideoFilter() *VideoFilter {
	return &VideoFilter{filterBase: newFilterBase()}
}

func (v *VideoFilter) SetIndex(i int) *VideoFilter { v.StreamIndex = i; return v }

func (v *VideoFilter) GetParameters() []string {
	return v.buildParameters("v")
}

func (v *VideoFilter) Scale(width, height int, flags ScalerFlags) *VideoFilter {
	v.addFilter(fmt.Sprintf("scale=%d:%d:flags=%s", width, height, flags))
	return v
}

func (v *VideoFilter) Crop(left, top, right, bottom int) *VideoFilter {
	v.addFilter(fmt.Sprintf("crop=w=in_w-%d-%d:h=in_h-%d-%d:x=%d:y=%d", left, right, top, bottom, left, top))
	return v
}

func (v *VideoFilter) Format(format ColorFormat) *VideoFilter {
	v.addFilter(fmt.Sprintf("format=pix_fmts=%s", format))
	return v
}

func (v *VideoFilter) HFlip() *VideoFilter { v.addFilter("hflip"); return v }
func (v *VideoFilter) VFlip() *VideoFilter { v.addFilter("vflip"); return v }

func (v *VideoFilter) Rotate(rotate string) *VideoFilter {
	v.addFilter(fmt.Sprintf("rotate=%s", rotate))
	return v
}

//////////
// AudioFilter
//////////

type AudioFilter struct {
	filterBase
}

func NewAudioFilter() *AudioFilter {
	return &AudioFilter{filterBase: newFilterBase()}
}

func (a *AudioFilter) SetIndex(i int) *AudioFilter { a.StreamIndex = i; return a }

func (a *AudioFilter) GetParameters() []string {
	return a.buildParameters("a")
}

func (a *AudioFilter) Volumedetect() *AudioFilter { a.addFilter("volumedetect"); return a }
func (a *AudioFilter) Replaygain() *AudioFilter   { a.addFilter("replaygain"); return a }

func (a *AudioFilter) Volume(volumeFilter string) *AudioFilter {
	a.addFilter(fmt.Sprintf("volume=%s", volumeFilter))
	return a
}

func (a *AudioFilter) LoudnormPrint() *AudioFilter {
	a.addFilter("loudnorm=print_format=json")
	return a
}

func (a *AudioFilter) Loudnorm(i, lra, tp float64, firstPass *LoudnormStats) *AudioFilter {
	s := fmt.Sprintf("loudnorm=I=%.1f:LRA=%.1f:TP=%.1f", i, lra, tp)
	if firstPass != nil {
		s += fmt.Sprintf(":measured_I=%.4f:measured_LRA=%.4f:measured_TP=%.4f:measured_thresh=%.4f",
			firstPass.I, firstPass.Lra, firstPass.TruePeak, firstPass.Threshold)
	}
	a.addFilter(s)
	return a
}

func (a *AudioFilter) Ebur128() *AudioFilter {
	a.addFilter("ebur128=metadata=1:peak=true")
	return a
}

func (a *AudioFilter) ChannelMap(channelMap string) *AudioFilter {
	a.addFilter(fmt.Sprintf("channelmap=%s", channelMap))
	return a
}

func (a *AudioFilter) Tempo(tempoMap string) *AudioFilter {
	a.addFilter(fmt.Sprintf("atempo=%s", tempoMap))
	return a
}

func (a *AudioFilter) TempoFpsToFps(fpsFrom, fpsTo string) *AudioFilter {
	a.addFilter(fmt.Sprintf("atempo=(%s)/(%s)", fpsTo, fpsFrom))
	return a
}

func (a *AudioFilter) TempoNtscToPal() *AudioFilter {
	return a.TempoFpsToFps("24000/1001", "25")
}

func (a *AudioFilter) TempoPalToNtsc() *AudioFilter {
	return a.TempoFpsToFps("25", "24000/1001")
}
