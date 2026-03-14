package goffmpeg

import (
	"fmt"
	"strings"
)

type metadataValue struct {
	Value        string
	PackInQuotes bool
}

type Metadata struct {
	StreamIndex int
	StreamType  *StreamType
	values      map[string]metadataValue
}

func NewMetadata(streamIndex int, streamType *StreamType) *Metadata {
	return &Metadata{
		StreamIndex: streamIndex,
		StreamType:  streamType,
		values:      make(map[string]metadataValue),
	}
}

func (m *Metadata) set(key, value string, pack bool) *Metadata {
	m.values[key] = metadataValue{Value: value, PackInQuotes: pack}
	return m
}

func (m *Metadata) Title(title string) *Metadata { return m.set("title", title, true) }
func (m *Metadata) TitleAppend(s string) *Metadata {
	v := m.values["title"]
	v.Value += s
	m.values["title"] = v
	return m
}
func (m *Metadata) Language(lang string) *Metadata { return m.set("language", lang, true) }
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
	return m.set("rotate", fmt.Sprintf("%d", degrees), false)
}
func (m *Metadata) Custom(key, value string, pack bool) *Metadata { return m.set(key, value, pack) }

func (m *Metadata) GetParameterStrings() []string {
	var prefix strings.Builder
	prefix.WriteString("-metadata")
	if m.StreamIndex >= 0 {
		prefix.WriteString(":s")
	}
	if m.StreamType != nil {
		prefix.WriteString(fmt.Sprintf(":%s", *m.StreamType))
	}
	if m.StreamIndex >= 0 {
		prefix.WriteString(fmt.Sprintf(":%d", m.StreamIndex))
	}
	p := prefix.String()

	var result []string
	for k, v := range m.values {
		val := v.Value
		if v.PackInQuotes {
			val = fmt.Sprintf("%q", val)
		}
		result = append(result, fmt.Sprintf("%s %s=%s", p, k, val))
	}
	return result
}
