package goffmpeg

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LoudnormStats struct {
	ExitCode  *int
	I         float64
	Threshold float64
	TruePeak  float64
	Lra       float64
}

func (s *LoudnormStats) String() string {
	if s.ExitCode != nil && *s.ExitCode != 0 {
		return fmt.Sprintf("Failed: %d", *s.ExitCode)
	}
	return fmt.Sprintf("I: %g LUFS, Threshold: %g LUFS, TruePeak: %g dBTP, Lra: %g LU",
		s.I, s.Threshold, s.TruePeak, s.Lra)
}

type VolumeStats struct {
	ExitCode *int
	Max      float64
	Mean     float64
}

func (s *VolumeStats) String() string {
	if s.ExitCode != nil && *s.ExitCode != 0 {
		return fmt.Sprintf("Failed: %d", *s.ExitCode)
	}
	return fmt.Sprintf("Mean: %g dB, Max: %g dB", s.Mean, s.Max)
}

type Ebur128Stats struct {
	ExitCode       *int
	I              float64
	Threshold      float64
	Peak           float64
	RangeLra       float64
	RangeThreshold float64
	RangeLraLow    float64
	RangeLraHigh   float64
}

func (s *Ebur128Stats) String() string {
	if s.ExitCode != nil && *s.ExitCode != 0 {
		return fmt.Sprintf("Failed: %d", *s.ExitCode)
	}
	return fmt.Sprintf("I: %g LUFS, Threshold: %g LUFS, Peak: %g dBFS, LRA: %g LU",
		s.I, s.Threshold, s.Peak, s.RangeLra)
}

type ReplaygainStats struct {
	ExitCode  *int
	TrackGain float64
	TrackPeak float64
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

// minFloat64 sentinel for "not found yet"
const notFound = -math.MaxFloat64

// parseFloat is a convenience wrapper around strconv.ParseFloat that uses
// the invariant (period) decimal separator, matching C#'s Convert.ToDouble.
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// runWithStderr runs an external command and calls lineHandler for every line
// written to stderr. It returns the process exit code and any execution error.
func runWithStderr(executable string, args []string, lineHandler func(string)) (int, error) {
	cmd := exec.Command(executable, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return -1, err
	}
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	scanLines(stderr, lineHandler)
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

// runWithStdout is like runWithStderr but captures stdout.
func runWithStdout(executable string, args []string, lineHandler func(string)) (int, error) {
	cmd := exec.Command(executable, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	scanLines(stdout, lineHandler)
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

func scanLines(r io.Reader, fn func(string)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fn(scanner.Text())
	}
}

// ─────────────────────────────────────────────
// Stats methods on Runner
// ─────────────────────────────────────────────

// GetReplaygainStats runs an FFmpeg job that emits replaygain metadata and
// parses track_gain / track_peak from stderr.
func (r *Runner) GetReplaygainStats(arguments *FFmpegArguments) *ReplaygainStats {
	trackGainRe := regexp.MustCompile(`.*track_gain = (.*) dB`)
	trackPeakRe := regexp.MustCompile(`.*track_peak = (.*)`)

	foundGain := notFound
	foundPeak := notFound

	argFields := strings.Fields(arguments.ArgumentString())
	exitCode, _ := runWithStderr(r.Settings.ExecutablePath, argFields, func(line string) {
		fmt.Println(line) // mirror C# Console.WriteLine
		if m := trackGainRe.FindStringSubmatch(line); m != nil {
			if v, err := parseFloat(m[1]); err == nil {
				foundGain = v
			}
		}
		if m := trackPeakRe.FindStringSubmatch(line); m != nil {
			if v, err := parseFloat(m[1]); err == nil {
				foundPeak = v
			}
		}
	})

	return &ReplaygainStats{
		ExitCode:  Ptr(exitCode),
		TrackGain: foundGain,
		TrackPeak: foundPeak,
	}
}

// GetLoudnormStatsFromArgs runs FFmpeg and parses loudnorm JSON output from stderr.
func (r *Runner) GetLoudnormStatsFromArgs(arguments *FFmpegArguments) *LoudnormStats {
	valueRe := regexp.MustCompile(`\s"(.+?)"\s:\s"([0-9\-\.]+)",?`)

	foundI := notFound
	foundThreshold := notFound
	foundTruePeak := notFound
	foundLra := notFound

	argFields := strings.Fields(arguments.ArgumentString())
	exitCode, _ := runWithStderr(r.Settings.ExecutablePath, argFields, func(line string) {
		m := valueRe.FindStringSubmatch(line)
		if m == nil {
			return
		}
		key, rawVal := m[1], m[2]
		v, err := parseFloat(rawVal)
		if err != nil {
			return
		}
		switch key {
		case "input_i":
			foundI = v
		case "input_tp":
			foundTruePeak = v
		case "input_lra":
			foundLra = v
		case "input_thresh":
			foundThreshold = v
		}
	})

	return &LoudnormStats{
		ExitCode:  Ptr(exitCode),
		I:         foundI,
		Threshold: foundThreshold,
		TruePeak:  foundTruePeak,
		Lra:       foundLra,
	}
}

// GetLoudnormStats is a convenience overload that builds the FFmpegArguments internally.
func (r *Runner) GetLoudnormStats(src string, index int) *LoudnormStats {
	st := StreamTypeAudio
	args := NewFFmpegArguments().
		SetLogLevel(LogLevelInfo).
		AddInputPath(src).
		SetNullOutput().
		Disable(true, false, true, true).
		AddMapping(NewMapping(0).SetStreamType(StreamTypeAudio)).
		AddAudioFilter(NewAudioFilter().LoudnormPrint().SetIndex(index))
	_ = st
	return r.GetLoudnormStatsFromArgs(args)
}

// GetVolumeStats runs volumedetect and parses mean/max volume from stderr.
func (r *Runner) GetVolumeStats(src string, index int) *VolumeStats {
	maxRe := regexp.MustCompile(`\[Parsed_volumedetect.*\].*max_volume:\s*([\-0-9\.]+)\s*dB`)
	meanRe := regexp.MustCompile(`\[Parsed_volumedetect.*\].*mean_volume:\s*([\-0-9\.]+)\s*dB`)

	foundMax := notFound
	foundMean := notFound

	args := NewFFmpegArguments().
		SetLogLevel(LogLevelInfo).
		AddInputPath(src).
		SetNullOutput().
		Disable(true, false, true, true).
		AddMapping(NewMapping(0).SetStreamType(StreamTypeAudio)).
		AddAudioFilter(NewAudioFilter().Volumedetect().SetIndex(index))

	argFields := strings.Fields(args.ArgumentString())
	exitCode, _ := runWithStderr(r.Settings.ExecutablePath, argFields, func(line string) {
		if m := meanRe.FindStringSubmatch(line); m != nil {
			if v, err := parseFloat(m[1]); err == nil {
				foundMean = v
			}
		}
		if m := maxRe.FindStringSubmatch(line); m != nil {
			if v, err := parseFloat(m[1]); err == nil {
				foundMax = v
			}
		}
	})

	return &VolumeStats{
		ExitCode: Ptr(exitCode),
		Max:      foundMax,
		Mean:     foundMean,
	}
}

// GetEbur128Stats runs the ebur128 filter and parses its summary from stderr.
func (r *Runner) GetEbur128Stats(src string, index int) *Ebur128Stats {
	parseRe := regexp.MustCompile(`\s+[a-zA-Z ]+:\s+([\-0-9\.]+)\s+.*`)

	foundI := notFound
	foundThreshold := notFound
	foundPeak := notFound
	foundRangeLra := notFound
	foundRangeThreshold := notFound
	foundRangeLraLow := notFound
	foundRangeLraHigh := notFound

	mode := 0

	args := NewFFmpegArguments().
		SetLogLevel(LogLevelInfo).
		AddInputPath(src).
		SetNullOutput().
		Disable(true, false, true, true).
		AddMapping(NewMapping(0).SetStreamType(StreamTypeAudio)).
		AddAudioFilter(NewAudioFilter().Ebur128().SetIndex(index))

	argFields := strings.Fields(args.ArgumentString())
	exitCode, _ := runWithStderr(r.Settings.ExecutablePath, argFields, func(line string) {
		extract := func() (float64, bool) {
			m := parseRe.FindStringSubmatch(line)
			if m == nil {
				return 0, false
			}
			v, err := parseFloat(m[1])
			return v, err == nil
		}

		switch {
		case strings.Contains(line, "Integrated loudness:"):
			mode = 1
		case mode == 1 && strings.Contains(line, "I:"):
			if v, ok := extract(); ok {
				foundI = v
			}
		case mode == 1 && strings.Contains(line, "Threshold:"):
			if v, ok := extract(); ok {
				foundThreshold = v
			}
			mode = 0
		case strings.Contains(line, "Loudness range:"):
			mode = 2
		case mode == 2 && strings.Contains(line, "LRA:"):
			if v, ok := extract(); ok {
				foundRangeLra = v
			}
		case mode == 2 && strings.Contains(line, "Threshold:"):
			if v, ok := extract(); ok {
				foundRangeThreshold = v
			}
		case mode == 2 && strings.Contains(line, "LRA low:"):
			if v, ok := extract(); ok {
				foundRangeLraLow = v
			}
		case mode == 2 && strings.Contains(line, "LRA high:"):
			if v, ok := extract(); ok {
				foundRangeLraHigh = v
			}
			mode = 0
		case strings.Contains(line, "True peak:"):
			mode = 3
		case mode == 3 && strings.Contains(line, "Peak:"):
			if v, ok := extract(); ok {
				foundPeak = v
			}
			mode = 0
		}
	})

	return &Ebur128Stats{
		ExitCode:       Ptr(exitCode),
		I:              foundI,
		Threshold:      foundThreshold,
		Peak:           foundPeak,
		RangeLra:       foundRangeLra,
		RangeThreshold: foundRangeThreshold,
		RangeLraLow:    foundRangeLraLow,
		RangeLraHigh:   foundRangeLraHigh,
	}
}

// ─────────────────────────────────────────────
// Probe helpers
// ─────────────────────────────────────────────

// GetKeyframes uses ffprobe to return a list of keyframe timestamps (in seconds).
func (r *Runner) GetKeyframes(src string) ([]float64, error) {
	var keyframes []float64
	args := []string{
		"-v", "error",
		"-skip_frame", "nokey",
		"-show_entries", "frame=pkt_pts_time",
		"-select_streams", "v",
		"-of", "csv=p=0",
		src,
	}
	_, err := runWithStdout(r.Settings.FfprobePath, args, func(line string) {
		line = strings.TrimSpace(line)
		if line == "" {
			return
		}
		if v, err := parseFloat(line); err == nil {
			keyframes = append(keyframes, v)
		}
	})
	return keyframes, err
}

// GetDuration uses ffprobe to return the duration of a media file.
func (r *Runner) GetDuration(src string) (time.Duration, error) {
	durationRe := regexp.MustCompile(`\s*Duration: (.*?),.*`)
	var result time.Duration

	args := []string{"-hide_banner", src}
	_, err := runWithStderr(r.Settings.FfprobePath, args, func(line string) {
		if !strings.Contains(line, "Duration") {
			return
		}
		m := durationRe.FindStringSubmatch(line)
		if m == nil {
			return
		}
		// Parse "hh:mm:ss.ff" — Go's time.Parse reference time is "15:04:05.99"
		t, err := time.Parse("15:04:05.00", strings.TrimSpace(m[1]))
		if err != nil {
			return
		}
		result = time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute +
			time.Duration(t.Second())*time.Second +
			time.Duration(t.Nanosecond())
	})
	return result, err
}

// PrintProbe prints Stream / Duration lines from ffprobe output to stdout.
func (r *Runner) PrintProbe(src string) error {
	args := []string{"-hide_banner", src}
	_, err := runWithStderr(r.Settings.FfprobePath, args, func(line string) {
		if strings.Contains(line, "Stream") ||
			strings.Contains(line, "DURATION") ||
			strings.Contains(line, "Duration") {
			fmt.Println(line)
		}
	})
	return err
}

// ─────────────────────────────────────────────
// Amplify optimiser
// ─────────────────────────────────────────────

// FindOptimalAmplifyValue iteratively adjusts an amplify value to bring the
// integrated loudness of videoFile as close as possible to targetI LUFS
// without exceeding maxPeak dBFS.
//
// setAction is called before each measurement so callers can rebuild their
// FFmpeg filter graph with the new amplify value.
func (r *Runner) FindOptimalAmplifyValue(
	initialAmplify float64,
	videoFile string,
	setAction func(float64),
	targetI float64,
	maxPeak float64,
	maxTries int,
) float64 {
	currAmplify := initialAmplify
	bestAmplify := initialAmplify
	bestDiff := 100.0
	tried := make(map[float64]bool)

	for {
		setAction(currAmplify)

		stats := r.GetEbur128Stats(videoFile, 0)
		fmt.Printf("Amplify %g -> Ebur128 - %s\n", currAmplify, stats)
		tried[currAmplify] = true

		iDiff := math.Abs(stats.I - targetI)
		if iDiff < bestDiff {
			if stats.Peak <= maxPeak {
				bestDiff = iDiff
				bestAmplify = currAmplify
			}
		} else {
			// Getting worse — stop.
			break
		}

		increment := 0.1
		if iDiff >= 5 {
			increment = 0.5
		} else if iDiff >= 2 {
			increment = 0.2
		}

		switch {
		case stats.Peak > maxPeak:
			currAmplify -= 0.1
		case stats.I > targetI:
			currAmplify -= increment
		case stats.I < targetI && stats.Peak <= maxPeak:
			currAmplify += increment
		}

		maxTries--
		if maxTries == 0 || tried[currAmplify] {
			break
		}
	}

	if bestAmplify == initialAmplify {
		fmt.Printf("Best amplify stays at %g\n", initialAmplify)
	} else {
		fmt.Printf("Best amplify should be %g (from initial %g)\n", bestAmplify, initialAmplify)
	}

	return bestAmplify
}
