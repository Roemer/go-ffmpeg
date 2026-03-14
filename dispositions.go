package goffmpeg

import (
	"fmt"
	"strings"
)

type Disposition struct {
	StreamIndex  int
	StreamType   StreamType
	dispositions []string
	isClear      bool
}

func NewDisposition(streamIndex int, streamType StreamType) *Disposition {
	return &Disposition{StreamIndex: streamIndex, StreamType: streamType}
}

func (d *Disposition) Clear() *Disposition { d.dispositions = nil; d.isClear = true; return d }
func (d *Disposition) Default() *Disposition {
	d.dispositions = append(d.dispositions, "default")
	return d
}
func (d *Disposition) Forced() *Disposition {
	d.dispositions = append(d.dispositions, "forced")
	return d
}
func (d *Disposition) Original() *Disposition {
	d.dispositions = append(d.dispositions, "original")
	return d
}
func (d *Disposition) Dub() *Disposition { d.dispositions = append(d.dispositions, "dub"); return d }

func (d *Disposition) GetParameters() []string {
	var value string
	if d.isClear {
		value = "0"
	} else {
		value = strings.Join(d.dispositions, "+")
	}
	if value == "" {
		return nil
	}
	return []string{fmt.Sprintf("-disposition:%s:%d", d.StreamType, d.StreamIndex), value}
}
