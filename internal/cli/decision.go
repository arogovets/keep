package cli

import "time"

type Decision string

const (
	DecisionMove Decision = "move"
	DecisionSkip Decision = "skip"
)

const syntheticActivityTolerance = 50 * time.Millisecond

func ShouldMove(always bool, idleFor time.Duration, interval time.Duration) Decision {
	return ShouldMoveAfterSynthetic(always, idleFor, interval, 0, false)
}

func ShouldMoveAfterSynthetic(always bool, idleFor time.Duration, interval time.Duration, syntheticAge time.Duration, hasSynthetic bool) Decision {
	if always || idleFor >= interval {
		return DecisionMove
	}

	if !hasSynthetic || syntheticAge+syntheticActivityTolerance < interval {
		return DecisionSkip
	}

	lastInputOffsetFromSynthetic := syntheticAge - idleFor
	if absDuration(lastInputOffsetFromSynthetic) <= syntheticActivityTolerance {
		return DecisionMove
	}
	return DecisionSkip
}

func absDuration(duration time.Duration) time.Duration {
	if duration < 0 {
		return -duration
	}
	return duration
}
