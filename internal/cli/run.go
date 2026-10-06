package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ieffai/stay-awake-cli/internal/activity"
	"github.com/ieffai/stay-awake-cli/internal/mover"
)

type lastActivityProvider interface {
	LastInputAge() (time.Duration, error)
	AccessibilityTrusted() bool
}

type cursorMover interface {
	MoveAndReturn(ctx context.Context, distance int) error
}

type actionReporter interface {
	StartupDetail(distance int) string
	ActionDescription() string
}

func Run(args []string, stdout io.Writer, stderr io.Writer, version string) (int, error) {
	cfg, err := ParseArgs(args, stderr, version)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0, nil
		}
		return 2, err
	}
	if cfg.Version {
		fmt.Fprintf(stdout, "keep %s\n", version)
		return 0, nil
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runner := appRunner{
		cfg:      cfg,
		out:      stdout,
		activity: activity.System{},
		mover:    mover.System{},
		now:      time.Now,
	}
	return runner.run(ctx)
}

type appRunner struct {
	cfg             Config
	out             io.Writer
	activity        lastActivityProvider
	mover           cursorMover
	now             func() time.Time
	lastSyntheticAt time.Time
}

func (r *appRunner) run(ctx context.Context) (int, error) {
	if !r.activity.AccessibilityTrusted() {
		return 1, errors.New(`Keep requires Accessibility permission.
Open System Settings -> Privacy & Security -> Accessibility
and allow your terminal application or the keep binary.`)
	}
	r.log("Keep started")
	fmt.Fprintf(r.out, "Interval: %s\n", formatDuration(r.cfg.Interval))
	if reporter, ok := r.mover.(actionReporter); ok {
		fmt.Fprintln(r.out, reporter.StartupDetail(r.cfg.Distance))
	} else {
		fmt.Fprintf(r.out, "Move distance: %d display point(s)\n", r.cfg.Distance)
	}
	fmt.Fprintln(r.out, "Press Ctrl+C to stop")
	fmt.Fprintln(r.out)

	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.log("stopped")
			return 0, nil
		case <-ticker.C:
			if err := r.tick(ctx); err != nil {
				r.log("error: %v", err)
			}
		}
	}
}

func (r *appRunner) tick(ctx context.Context) error {
	idleFor, err := r.activity.LastInputAge()
	if err != nil {
		return fmt.Errorf("read activity: %w", err)
	}

	var syntheticAge time.Duration
	hasSynthetic := !r.lastSyntheticAt.IsZero()
	now := r.now()
	if hasSynthetic {
		syntheticAge = now.Sub(r.lastSyntheticAt)
	}

	if ShouldMoveAfterSynthetic(r.cfg.Always, idleFor, r.cfg.Interval, syntheticAge, hasSynthetic) == DecisionSkip {
		if r.cfg.Verbose {
			r.log("user activity detected after %s idle - skipped", formatDuration(idleFor))
			return nil
		}
		r.log("user activity detected - skipped")
		return nil
	}

	if err := r.mover.MoveAndReturn(ctx, r.cfg.Distance); err != nil {
		return fmt.Errorf("move cursor: %w", err)
	}
	r.lastSyntheticAt = r.now()
	if r.cfg.Always && idleFor < r.cfg.Interval {
		r.log("always mode - %s after %s idle", r.actionDescription(), formatDuration(idleFor))
		return nil
	}
	r.log("idle for %s - %s", formatDuration(idleFor), r.actionDescription())
	return nil
}

func (r *appRunner) actionDescription() string {
	if reporter, ok := r.mover.(actionReporter); ok {
		return reporter.ActionDescription()
	}
	return "cursor moved"
}

func (r *appRunner) log(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(r.out, "%s %s\n", r.now().Format("15:04:05"), message)
}

func formatDuration(duration time.Duration) string {
	if duration%time.Second == 0 {
		return fmt.Sprintf("%ds", int(duration/time.Second))
	}
	return duration.String()
}
