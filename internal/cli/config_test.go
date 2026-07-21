package cli

import (
	"io"
	"testing"
	"time"
)

func TestParseArgsDefaults(t *testing.T) {
	cfg, err := ParseArgs(nil, io.Discard, "test")
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	if cfg.Interval != 60*time.Second {
		t.Fatalf("Interval = %s, want 60s", cfg.Interval)
	}
	if cfg.Distance != 1 {
		t.Fatalf("Distance = %d, want 1", cfg.Distance)
	}
	if cfg.Always {
		t.Fatal("Always = true, want false")
	}
	if cfg.Verbose {
		t.Fatal("Verbose = true, want false")
	}
}

func TestParseArgsOptions(t *testing.T) {
	cfg, err := ParseArgs([]string{"--interval", "2m", "--distance", "3", "--always", "--verbose"}, io.Discard, "test")
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	if cfg.Interval != 2*time.Minute {
		t.Fatalf("Interval = %s, want 2m", cfg.Interval)
	}
	if cfg.Distance != 3 {
		t.Fatalf("Distance = %d, want 3", cfg.Distance)
	}
	if !cfg.Always {
		t.Fatal("Always = false, want true")
	}
	if !cfg.Verbose {
		t.Fatal("Verbose = false, want true")
	}
}

func TestParseIntervalBareSeconds(t *testing.T) {
	got, err := ParseInterval("60")
	if err != nil {
		t.Fatalf("ParseInterval() error = %v", err)
	}
	if got != time.Minute {
		t.Fatalf("ParseInterval() = %s, want 1m", got)
	}
}

func TestParseArgsRejectsInvalidValues(t *testing.T) {
	cases := [][]string{
		{"--interval", "0s"},
		{"--interval", "soon"},
		{"--distance", "0"},
		{"extra"},
	}
	for _, args := range cases {
		if _, err := ParseArgs(args, io.Discard, "test"); err == nil {
			t.Fatalf("ParseArgs(%v) error = nil, want error", args)
		}
	}
}
