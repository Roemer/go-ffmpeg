package goffmpeg

import (
	"fmt"
	"strings"
)

type metadataValue struct {
	Value string
}

type Metadata struct {
	streamIndex int
	streamType  StreamType
	values      map[string]metadataValue
}

func NewMetadata() *Metadata {
	return &Metadata{
		streamIndex: 0,
		streamType:  StreamTypeNone,
		values:      make(map[string]metadataValue),
	}
}

func (m *Metadata) set(key, value string) *Metadata {
	m.values[key] = metadataValue{Value: value}
	return m
}

func (m *Metadata) StreamType(streamType StreamType) *Metadata {
	m.streamType = streamType
	return m
}

func (m *Metadata) StreamIndex(index int) *Metadata {
	m.streamIndex = index
	return m
}

func (m *Metadata) Title(title string) *Metadata { return m.set("title", title) }
func (m *Metadata) TitleAppend(s string) *Metadata {
	v := m.values["title"]
	v.Value += s
	m.values["title"] = v
	return m
}
func (m *Metadata) Language(lang string) *Metadata { return m.set("language", lang) }
func (m *Metadata) German() *Metadata              { m.Title("Deutsch"); return m.Language("deu") }
func (m *Metadata) English() *Metadata             { m.Title("English"); return m.Language("eng") }
func (m *Metadata) EnglishCommentary() *Metadata {
	m.Title("English (Commentary)")
	return m.Language("eng")
}
func (m *Metadata) Portuguese() *Metadata { m.Title("Português"); return m.Language("por") }
func (m *Metadata) NorvegianBokmal() *Metadata {
	m.Title("Norvegian Bokmål")
	return m.Language("nob")
}
func (m *Metadata) Korean() *Metadata      { m.Title("Korean"); return m.Language("kor") }
func (m *Metadata) French() *Metadata      { m.Title("Français"); return m.Language("fra") }
func (m *Metadata) Italian() *Metadata     { m.Title("Italiano"); return m.Language("ita") }
func (m *Metadata) SwissGerman() *Metadata { m.Title("Schweizerdeutsch"); return m.Language("gsw") }
func (m *Metadata) Japanese() *Metadata    { m.Title("Japanese"); return m.Language("jpn") }
func (m *Metadata) Rotate(degrees int) *Metadata {
	return m.set("rotate", fmt.Sprintf("%d", degrees))
}
func (m *Metadata) Custom(key, value string) *Metadata { return m.set(key, value) }

func (m *Metadata) GetParameters() []string {
	var prefix strings.Builder
	prefix.WriteString("-metadata")
	if m.streamIndex >= 0 {
		prefix.WriteString(":s")
	}
	if m.streamType != StreamTypeNone {
		fmt.Fprintf(&prefix, ":%s", m.streamType)
	}
	if m.streamIndex >= 0 {
		fmt.Fprintf(&prefix, ":%d", m.streamIndex)
	}
	p := prefix.String()

	var result []string
	for k, v := range m.values {
		val := v.Value
		result = append(result, p, fmt.Sprintf("%s=%s", k, val))
	}
	return result
}
