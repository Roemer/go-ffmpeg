package goffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	Settings *FFmpegSettings
}

func NewRunner(settings *FFmpegSettings) *Runner {
	if settings == nil {
		settings = DefaultFFmpegSettings()
	}
	return &Runner{Settings: settings}
}

func (r *Runner) ExecuteFFmpegArgs(arguments *FFmpegArguments) (int, error) {
	args := arguments.ArgumentSlice()
	if r.Settings.ShowCommandInConsole {
		fmt.Println(arguments.ArgumentString())
	}
	return r.executeCmd(r.Settings.ExecutablePath, args)
}

func (r *Runner) ExecuteFFmpegRaw(rawArgs ...string) (int, error) {
	settings := r.Settings
	args := []string{
		"-hide_banner",
		"-xerror",
		fmt.Sprintf("-loglevel %s", settings.LogLevel),
		"-stats",
		fmt.Sprintf("-stats_period %d", settings.StatsPeriod),
	}
	if settings.Overwrite {
		args = append(args, "-y")
	}
	args = append(args, rawArgs...)
	if settings.ShowCommandInConsole {
		fmt.Println(strings.Join(args, " "))
	}
	return r.executeCmd(settings.ExecutablePath, args)
}

func (r *Runner) executeCmd(executable string, args []string) (int, error) {
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if r.Settings.LogMessageAction != nil {
		// stderr forwarding via LogMessageAction would require more complex setup;
		// left as an extension point.
	}
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), err
		}
		return -1, err
	}
	return 0, nil
}
