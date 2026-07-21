package cli

import (
	"testing"
	"time"
)

func TestShouldMove(t *testing.T) {
	cases := []struct {
		name     string
		always   bool
		idleFor  time.Duration
		interval time.Duration
		want     Decision
	}{
		{name: "idle below interval skips", idleFor: 59 * time.Second, interval: time.Minute, want: DecisionSkip},
		{name: "idle at interval moves", idleFor: time.Minute, interval: time.Minute, want: DecisionMove},
		{name: "idle above interval moves", idleFor: 61 * time.Second, interval: time.Minute, want: DecisionMove},
		{name: "always moves regardless of idle", always: true, idleFor: time.Second, interval: time.Minute, want: DecisionMove},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldMove(tc.always, tc.idleFor, tc.interval)
			if got != tc.want {
				t.Fatalf("ShouldMove() = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestShouldMoveAfterSynthetic(t *testing.T) {
	cases := []struct {
		name         string
		idleFor      time.Duration
		syntheticAge time.Duration
		want         Decision
	}{
		{
			name:         "synthetic input at interval moves",
			idleFor:      time.Minute,
			syntheticAge: time.Minute,
			want:         DecisionMove,
		},
		{
			name:         "user input shortly after synthetic skips",
			idleFor:      59*time.Second + 900*time.Millisecond,
			syntheticAge: time.Minute,
			want:         DecisionSkip,
		},
		{
			name:         "too soon after synthetic skips",
			idleFor:      10 * time.Second,
			syntheticAge: 10 * time.Second,
			want:         DecisionSkip,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldMoveAfterSynthetic(false, tc.idleFor, time.Minute, tc.syntheticAge, true)
			if got != tc.want {
				t.Fatalf("ShouldMoveAfterSynthetic() = %s, want %s", got, tc.want)
			}
		})
	}
}
