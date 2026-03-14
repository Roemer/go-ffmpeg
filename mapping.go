package goffmpeg

import (
	"fmt"
	"strings"
)

type Mapping struct {
	Index       int
	StreamType  *StreamType
	StreamIndex int
	Negative    bool
	Optional    bool
}

func NewMapping(index int) *Mapping {
	return &Mapping{Index: index, StreamIndex: -1}
}

func (m *Mapping) SetStreamType(st StreamType) *Mapping {
	m.StreamType = &st
	return m
}

func (m *Mapping) SetStreamIndex(idx int) *Mapping {
	m.StreamIndex = idx
	return m
}

func (m *Mapping) SetNegative(negative bool) *Mapping {
	m.Negative = negative
	return m
}

func (m *Mapping) SetOptional(optional bool) *Mapping {
	m.Optional = optional
	return m
}

func (m *Mapping) GetParameterString() string {
	var sb strings.Builder
	sb.WriteString("-map ")
	if m.Negative {
		sb.WriteString("-")
	}
	fmt.Fprintf(&sb, "%d", m.Index)
	if m.StreamType != nil {
		fmt.Fprintf(&sb, ":%s", *m.StreamType)
	}
	if m.StreamIndex >= 0 {
		fmt.Fprintf(&sb, ":%d", m.StreamIndex)
	}
	if m.Optional {
		sb.WriteString("?")
	}
	return sb.String()
}
