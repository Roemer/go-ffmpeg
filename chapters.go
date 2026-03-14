package goffmpeg

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ChapterEntry struct {
	Time time.Duration
	Text string
}

func NewChapterEntry(t time.Duration, text string) *ChapterEntry {
	return &ChapterEntry{Time: t, Text: text}
}

func NewChapterEntryFromString(timeString, text string) (*ChapterEntry, error) {
	d, err := parseSRTDuration(timeString)
	if err != nil {
		return nil, fmt.Errorf("NewChapterEntryFromString: %w", err)
	}
	return NewChapterEntry(d, text), nil
}

func parseSRTDuration(s string) (time.Duration, error) {
	// Replace comma with dot so time.Parse can handle milliseconds.
	s = strings.ReplaceAll(s, ",", ".")
	t, err := time.Parse("15:04:05.000", s)
	if err != nil {
		return 0, fmt.Errorf("parseSRTDuration: invalid format %q: %w", s, err)
	}
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond()), nil
}

func CreateMetadataFromChapters(chapters []*ChapterEntry, outputPath string) (bool, error) {
	var sb strings.Builder
	sb.WriteString(";FFMETADATA1\n")

	for i, chapter := range chapters {
		sb.WriteString("\n")
		sb.WriteString("[CHAPTER]\n")
		sb.WriteString("TIMEBASE=1/1000\n")
		fmt.Fprintf(&sb, "START=%d\n", chapter.Time.Milliseconds())

		if i < len(chapters)-1 {
			next := chapters[i+1]
			fmt.Fprintf(&sb, "END=%d\n", next.Time.Milliseconds()-1)
		} else {
			// Last chapter — no known end time.
			sb.WriteString("END=\n")
		}

		if strings.TrimSpace(chapter.Text) != "" {
			sb.WriteString(fmt.Sprintf("title=%s\n", chapter.Text))
		}
	}

	if err := os.WriteFile(outputPath, []byte(sb.String()), 0o644); err != nil {
		return false, err
	}

	return len(chapters) > 0, nil
}
