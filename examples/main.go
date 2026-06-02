package main

import (
	"fmt"

	goffmpeg "github.com/roemer/go-ffmpeg"
)

func main() {
	args := goffmpeg.NewFFmpegArguments().
		AddInputPath(`JWCC - Hidden Adventure unedited version.mkv`).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeVideo)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeAudio).SetStreamIndex(4)).
		SetOutputPath(`adv_german.mp4`)
	fmt.Println(args.ArgumentString())
	return
	ExtractAudio("input.mkv", "output.mp3", 0)
	ExtractVideo("input.mkv", "output.mp4", 0)
	ConvertAudio("input.mkv", "output.aac", 0)
	ConvertVideo("input.mkv", "output.mp4", 1, &goffmpeg.X264Settings{
		Preset: goffmpeg.X264PresetVeryslow,
		Tune:   goffmpeg.X264TuneAnimation,
		CRF:    18,
	})

	/*args := goffmpeg.NewFFmpegArguments().
		AddInputPath(`\\nas03\Dump\Beavis And Butthead Do America crf20.mkv`).
		AddMapping(goffmpeg.NewMapping(1).SetStreamType(goffmpeg.StreamTypeVideo)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeAudio)).
		SetMapChaptersIndex(0).
		SetMaxInterleaveDelta(0)
	fmt.Println(args.ArgumentString())*/

}

func ExtractAudio(src, dst string, index int) {
	args := goffmpeg.NewFFmpegArguments().
		AddInputPath(src).
		Disable(true, false, true, true).
		SetCopyAll(true).
		SetOutput(goffmpeg.NewOutputFile(dst)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeAudio).SetStreamIndex(index))
	fmt.Println(args.ArgumentString())
}

func ExtractVideo(src, dst string, index int) {
	args := goffmpeg.NewFFmpegArguments().
		AddInputPath(src).
		Disable(false, true, true, true).
		SetCopyAll(true).
		SetOutput(goffmpeg.NewOutputFile(dst)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeVideo).SetStreamIndex(index))
	fmt.Println(args.ArgumentString())
}

func ConvertAudio(src, dst string, index int) {
	args := goffmpeg.NewFFmpegArguments().
		AddInputPath(src).
		Disable(true, false, true, true).
		AddAudioEncoding(goffmpeg.NewAacAudioEncoding()).
		SetOutput(goffmpeg.NewOutputFile(dst)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeAudio).SetStreamIndex(index))
	fmt.Println(args.ArgumentString())
}

func ConvertVideo(src, dst string, index int, x264Settings *goffmpeg.X264Settings) {
	videoEncoding := goffmpeg.NewX264VideoEncoding()
	if x264Settings != nil {
		videoEncoding.FromSettings(x264Settings)
	}
	args := goffmpeg.NewFFmpegArguments().
		AddInputPath(src).
		Disable(false, true, true, true).
		AddVideoEncoding(videoEncoding).
		SetOutput(goffmpeg.NewOutputFile(dst)).
		AddMapping(goffmpeg.NewMapping(0).SetStreamType(goffmpeg.StreamTypeVideo).SetStreamIndex(index))
	fmt.Println(args.ArgumentString())
}

/*
func (r *Runner) ChangeVideoPlayRateFps(inputFile, outputFile string, newFps float64) error {
	tempFile := os.TempDir() + "/ffmpeg_tmp.h264"
	defer os.Remove(tempFile)

	if _, err := r.ExecuteFFmpegRaw(
		fmt.Sprintf(`-i %q`, inputFile),
		"-map 0:v -vcodec copy -bsf:v h264_mp4toannexb",
		fmt.Sprintf("%q", tempFile),
	); err != nil {
		return err
	}
	_, err := r.ExecuteFFmpegRaw(
		"-fflags +genpts",
		fmt.Sprintf("-r %g", newFps),
		fmt.Sprintf(`-i %q`, tempFile),
		"-vcodec copy",
		fmt.Sprintf("%q", outputFile),
	)
	return err
}*/
