# stay

`stay` is a tiny macOS CLI utility that keeps the user session active by posting real mouse movement events through Core Graphics.

It is designed as a transparent foreground process:

```bash
stay
```

While the command is running, `stay` is active. Press `Ctrl+C` and it stops immediately.

`stay` is CLI-only. It does not install a `.app`, menu bar item, tray icon, daemon, launch agent, AppleScript, or any GUI component.

## Why

Some GUI mouse mover apps are hard to inspect, hard to debug, and awkward to run in controlled environments. `stay` keeps the behavior explicit:

- one terminal command starts it;
- terminal logs show every decision;
- `Ctrl+C`, `SIGINT`, or `SIGTERM` stops it;
- the cursor returns to its original position after every synthetic movement;
- there is no accumulated cursor drift.

The implementation is inspired by the general idea of periodic mouse movement tools, including `automatic-mouse-mover`, but no source code from that project is copied here.

## Requirements

- macOS.
- Apple Silicon or Intel Mac.
- Go 1.22+ if building from source.
- Accessibility permission for either:
  - the terminal application that runs `stay`; or
  - the installed `stay` binary.

`stay` uses macOS ApplicationServices/Core Graphics APIs. It does not use AppleScript.

## Install

Install with Homebrew:

```bash
brew install KarlinskyS/tap/stay
```

Then run:

```bash
stay
```

Homebrew builds `stay` from source and installs any build dependencies it needs.

## Build From Source

Build and install into `$HOME/.local/bin`:

```bash
make build
make install
```

Make sure `$HOME/.local/bin` is in your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Install into `/usr/local/bin` instead:

```bash
sudo make install PREFIX=/usr/local
```

## Run

```bash
stay
```

Default behavior:

- checks activity every `60s`;
- moves the cursor by `1px`;
- skips synthetic movement if real mouse or keyboard activity happened during the interval;
- moves to a nearby valid on-screen point and then returns to the original cursor position.

Example output:

```text
13:42:00 Stay started
Interval: 60s
Move distance: 1px
Press Ctrl+C to stop

13:43:00 idle for 60s - cursor moved
13:44:00 user activity detected - skipped
13:45:00 idle for 60s - cursor moved
13:45:12 stopped
```

## Options

```bash
stay --interval 60s
stay --distance 1
stay --always
stay --verbose
stay --version
stay --help
```

| Option | Description |
| --- | --- |
| `--interval` | Check interval. Accepts Go duration values such as `60s`, `1m`, `500ms`. A bare integer is treated as seconds. |
| `--distance` | Cursor movement distance in pixels. Default: `1`. |
| `--always` | Move every interval, even if user activity was detected. |
| `--verbose` | Print more detailed diagnostics. |
| `--version` | Print version and exit. |
| `--help` | Print usage and exit. |

Examples:

```bash
stay --interval 30s
stay --interval 5s --verbose
stay --distance 2
stay --always
```

## Accessibility Permission

On first run, macOS may block synthetic input events until Accessibility permission is granted.

If permission is missing, `stay` exits with:

```text
Stay requires Accessibility permission.
Open System Settings -> Privacy & Security -> Accessibility
and allow your terminal application or the stay binary.
```

To grant permission:

1. Open `System Settings`.
2. Go to `Privacy & Security`.
3. Open `Accessibility`.
4. Enable your terminal application, for example `Terminal`, `iTerm2`, `Ghostty`, `WezTerm`, or `Alacritty`.
5. Run `stay` again.

If you run the installed binary directly and macOS shows the binary as the requesting process, allow `stay` itself.

## How It Works

Every interval, `stay` asks macOS how long it has been since the last input event.

If the user has been idle for at least the configured interval, `stay`:

1. reads the current cursor position;
2. reads the active display bounds;
3. chooses a nearby valid point on any active display;
4. posts a Core Graphics `kCGEventMouseMoved` event;
5. waits briefly;
6. posts another mouse move event back to the original position.

This handles screen edges and multi-monitor layouts, including displays with negative coordinates.

Because synthetic movement itself resets the macOS idle timer, `stay` tracks its own last synthetic movement and distinguishes it from real user activity.

## Diagnostics

Run with a short interval while testing:

```bash
stay --interval 5s --verbose
```

Check macOS power assertions:

```bash
pmset -g assertions
```

After synthetic movement, macOS should report recent user activity, commonly through a `UserIsActive` assertion owned by `WindowServer`.

Useful checks:

```bash
stay --version
stay --help
make test
make build
```

## Troubleshooting

### `command not found: stay`

The install path is probably not in `PATH`.

Run directly:

```bash
$HOME/.local/bin/stay
```

Or add it to your shell profile:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Accessibility permission was granted, but movement still fails

Try the following:

- restart the terminal application;
- remove and re-add the terminal application in Accessibility settings;
- if launching `/usr/local/bin/stay` or `$HOME/.local/bin/stay` directly, allow the binary if macOS lists it separately;
- run `stay --interval 5s --verbose` and inspect the logs.

### Cursor moves but a remote desktop session does not see it

Remote desktop clients decide which local input events are forwarded into the remote session. `stay` can only generate local macOS mouse movement. The remote client must accept and forward that movement.

## Uninstall

If installed into `$HOME/.local/bin`:

```bash
rm -f "$HOME/.local/bin/stay"
```

If installed into `/usr/local/bin`:

```bash
sudo rm -f /usr/local/bin/stay
```

You can also remove the Accessibility permission from `System Settings -> Privacy & Security -> Accessibility`.

## Development

```bash
make test
make build
bin/stay --help
bin/stay --interval 5s --verbose
```

Project layout:

```text
.
├── cmd/stay/main.go
├── internal/activity
├── internal/cli
├── internal/mover
├── go.mod
├── Makefile
├── README.md
└── LICENSE
```

Manual smoke checklist:

- `stay` starts a long-running foreground process.
- Terminal output clearly shows that the process is active.
- `Ctrl+C` stops the process and logs shutdown.
- `SIGTERM` stops the process and logs shutdown.
- Each actual synthetic movement is logged.
- Cursor returns to the original position after movement.
- Repeated movement does not accumulate cursor drift.
- Movement works on the primary monitor.
- Movement works on an additional monitor.
- Cursor near screen edges still moves to a valid on-screen point and returns.
- Missing Accessibility permission shows the documented error.
- `pmset -g assertions` or another system indicator sees user activity after synthetic movement.

## Policy

Use `stay` only in ways that comply with your organization policies, contracts, and local rules. This tool is intentionally transparent so its behavior is easy to inspect and reason about.

## License

MIT.
