package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultInterval = 60 * time.Second
	DefaultDistance = 1
)

type Config struct {
	Interval time.Duration
	Distance int
	Always   bool
	Verbose  bool
	Version  bool
}

func DefaultConfig() Config {
	return Config{
		Interval: DefaultInterval,
		Distance: DefaultDistance,
	}
}

func ParseArgs(args []string, out io.Writer, version string) (Config, error) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("keep", flag.ContinueOnError)
	fs.SetOutput(out)

	interval := fs.String("interval", cfg.Interval.String(), "check interval, for example 60s, 1m, or 500ms")
	fs.IntVar(&cfg.Distance, "distance", cfg.Distance, "cursor movement distance in display points")
	fs.BoolVar(&cfg.Always, "always", cfg.Always, "move cursor every interval regardless of user activity")
	fs.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "print detailed diagnostic logs")
	fs.BoolVar(&cfg.Version, "version", false, "print version and exit")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: keep [options]\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if fs.NArg() > 0 {
		return cfg, fmt.Errorf("unexpected argument: %s", fs.Arg(0))
	}

	parsedInterval, err := ParseInterval(*interval)
	if err != nil {
		return cfg, err
	}
	cfg.Interval = parsedInterval
	if cfg.Distance < 1 {
		return cfg, errors.New("--distance must be greater than 0")
	}

	return cfg, nil
}

func ParseInterval(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("--interval must not be empty")
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds <= 0 {
			return 0, errors.New("--interval must be greater than 0")
		}
		return time.Duration(seconds) * time.Second, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid --interval %q: use values like 60s, 1m, or 500ms", value)
	}
	if duration <= 0 {
		return 0, errors.New("--interval must be greater than 0")
	}
	return duration, nil
}
