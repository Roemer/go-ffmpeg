package goffmpeg

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type SubtitleEntry struct {
	TimeString string
	Text       string
}

func NewSubtitleEntryRaw(timeString, text string) *SubtitleEntry {
	return &SubtitleEntry{TimeString: timeString, Text: text}
}

func NewSubtitleEntry(from, to time.Duration, text string) *SubtitleEntry {
	return &SubtitleEntry{
		TimeString: fmt.Sprintf("%s --> %s", formatSRTDuration(from), formatSRTDuration(to)),
		Text:       text,
	}
}

func formatSRTDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}

func GenerateSubtitleFile(outputFilePath string, entries []*SubtitleEntry) error {
	if len(entries) == 0 {
		// Best-effort delete; ignore "file not found" errors.
		if err := os.Remove(outputFilePath); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	var sb strings.Builder
	for i, sub := range entries {
		fmt.Fprintf(&sb, "%d\n", i+1)
		fmt.Fprintf(&sb, "%s\n", sub.TimeString)
		fmt.Fprintf(&sb, "%s\n", sub.Text)
		fmt.Fprintf(&sb, "\n")
	}

	return os.WriteFile(outputFilePath, []byte(sb.String()), 0o644)
}
