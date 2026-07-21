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
	got := ShouldMoveAfterSynthetic(false, 59*time.Second+900*time.Millisecond, time.Minute, time.Minute, true)
	if got != DecisionMove {
		t.Fatalf("ShouldMoveAfterSynthetic() = %s, want %s", got, DecisionMove)
	}

	got = ShouldMoveAfterSynthetic(false, 10*time.Second, time.Minute, time.Minute, true)
	if got != DecisionSkip {
		t.Fatalf("ShouldMoveAfterSynthetic() = %s, want %s", got, DecisionSkip)
	}
}
