package cli

import "time"

type Decision string

const (
	DecisionMove Decision = "move"
	DecisionSkip Decision = "skip"
)

const syntheticActivityTolerance = 250 * time.Millisecond

func ShouldMove(always bool, idleFor time.Duration, interval time.Duration) Decision {
	return ShouldMoveAfterSynthetic(always, idleFor, interval, 0, false)
}

func ShouldMoveAfterSynthetic(always bool, idleFor time.Duration, interval time.Duration, syntheticAge time.Duration, hasSynthetic bool) Decision {
	if always || idleFor >= interval {
		return DecisionMove
	}
	if hasSynthetic &&
		syntheticAge+syntheticActivityTolerance >= interval &&
		idleFor+syntheticActivityTolerance >= syntheticAge {
		return DecisionMove
	}
	return DecisionSkip
}
