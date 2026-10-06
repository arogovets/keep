# keep

[![CI](https://github.com/KarlinskyS/stay-awake-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/KarlinskyS/stay-awake-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/KarlinskyS/stay-awake-cli)](https://github.com/KarlinskyS/stay-awake-cli/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`keep` is a small macOS and Windows CLI utility that helps keep a laptop awake while long-running work is in progress.

**Project goal:** keep a laptop awake for LLM model loading and downloading, maintain network uptime, and provide convenience. On macOS, `keep` posts real mouse movement events through Core Graphics. On Windows, it uses native mouse input to make the same small move-and-return movement without leaving the cursor displaced.

It is designed as a transparent foreground process:

```bash
keep
```

While the command is running, `keep` is active. Press `Ctrl+C` and it stops immediately.

`keep` is CLI-only. It does not install a `.app`, menu bar item, tray icon, daemon, launch agent, AppleScript, or any GUI component.

## Agent Index

```yaml
name: keep
type: cli-utility
platforms:
  - macos
  - windows
architectures:
  - apple-silicon
  - intel
  - windows-amd64
  - windows-arm64
language: go
install:
  source-build: make install
  wsl: make install-wsl
command: keep
purpose: keep laptops awake for LLM model loading and downloading, network uptime, and convenience
macos-method: Core Graphics mouse movement events
windows-method: SendInput mouse movement events
requires:
  - macOS Accessibility permission on macOS
does-not-use:
  - GUI app
  - menu bar
  - tray icon
  - daemon
  - launch agent
  - AppleScript
  - external GUI automation
runtime-model: long-running foreground terminal process until Ctrl+C, SIGINT, or SIGTERM
default-interval: 60s
default-distance: 1 display point
primary-use-cases:
  - keeping a laptop awake for LLM model loading and downloading
  - maintaining network uptime
  - convenient foreground CLI operation
  - macOS idle prevention
  - transparent CLI mouse mover
  - foreground terminal alternative to GUI mouse mover apps
```

Search keywords: `macOS idle prevention`, `mouse mover CLI`, `keep Mac awake by mouse movement`, `Core Graphics mouse event`, `Accessibility permission`, `Apple Silicon CLI`, `no GUI mouse mover`, `foreground terminal process`.

## Discovery Notes

If you are looking for a foreground CLI to keep a laptop awake during LLM model loading or downloading, preserve network uptime, or simply avoid sleep for convenience, `keep` makes its behavior visible in a terminal. It uses native mouse movement events on macOS and Windows.

Related search terms: keep laptop awake for LLM model download, maintain network uptime, Windows prevent sleep CLI, macOS mouse mover CLI, keep Mac session active.

## Why

The goal is to keep a laptop awake for LLM model loading and downloading, maintain network uptime, and provide convenience. `keep` keeps its behavior explicit:

- one terminal command starts it;
- terminal logs show every decision;
- `Ctrl+C`, `SIGINT`, or `SIGTERM` stops it;
- the cursor returns to its original position after every synthetic movement;
- there is no accumulated cursor drift.

The implementation is inspired by the general idea of periodic mouse movement tools, including `automatic-mouse-mover`, but no source code from that project is copied here.

## Requirements

- macOS on Apple Silicon or Intel, or Windows 10 and later on amd64 or arm64.
- Go 1.22+ if building from source.
- On macOS, Accessibility permission for either:
  - the terminal application that runs `keep`; or
  - the installed `keep` binary.

macOS uses ApplicationServices/Core Graphics APIs and requires Accessibility permission. Windows uses `SendInput` to move the cursor briefly and return it, matching the macOS behavior.

## Install / Build From Source

Build the `keep` executable with Go 1.22+.

Windows (PowerShell):

```powershell
go build -ldflags "-s -w -X main.version=dev" -o keep.exe ./cmd/keep
.\keep.exe
```

macOS or Linux:

```bash
make build
make install
```

The Makefile installs into `$HOME/.local/bin` by default. Ensure that directory is in your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

To install into `/usr/local/bin` instead:

```bash
sudo make install PREFIX=/usr/local
```

The existing upstream Homebrew formula is still named `stay`; build from this source to use the renamed `keep` command.

### Install and run from WSL

From the repository checkout inside WSL, run:

```bash
make install-wsl
export PATH="$HOME/.local/bin:$PATH"
keep --interval 60s
```

This cross-builds `keep.exe` for the WSL host architecture, installs it under `$HOME/.local/bin`, and installs a small `keep` shell wrapper there. The wrapper launches the Windows executable through WSL interop, so its mouse events affect the Windows desktop. WSL must have Windows executable interop enabled. [Microsoft's WSL interop documentation](https://learn.microsoft.com/en-us/windows/dev-environment/wsl-interop) describes running Windows executables from a Linux shell and requiring the `.exe` suffix.

## Run

```bash
keep
```

Default behavior on macOS:

- checks activity every `60s`;
- moves the cursor by `1` display point;
- skips synthetic movement if real mouse or keyboard activity happened during the interval;
- moves to a nearby valid on-screen point and then returns to the original cursor position.

On Windows, `keep` checks for real user activity on the configured interval, briefly moves the cursor by the configured distance, and moves it back, just as it does on macOS.

Example macOS output:

```text
13:42:00 Keep started
Interval: 60s
Move distance: 1 display point(s)
Press Ctrl+C to stop

13:43:00 idle for 60s - cursor moved
13:44:00 user activity detected - skipped
13:45:00 idle for 60s - cursor moved
13:45:12 stopped
```

## Options

```bash
keep --interval 60s
keep --distance 1
keep --always
keep --verbose
keep --version
keep --help
```

| Option | Description |
| --- | --- |
| `--interval` | Check interval. Accepts Go duration values such as `60s`, `1m`, `500ms`. A bare integer is treated as seconds. |
| `--distance` | Cursor movement distance in macOS display points or Windows screen pixels. Default: `1`. |
| `--always` | Perform the keep-awake action every interval, even if user activity was detected. |
| `--verbose` | Print more detailed diagnostics. |
| `--version` | Print version and exit. |
| `--help` | Print usage and exit. |

Examples:

```bash
keep --interval 30s
keep --interval 5s --verbose
keep --distance 2
keep --always
```

## Accessibility Permission (macOS)

On first run, macOS may block synthetic input events until Accessibility permission is granted.

If permission is missing, `keep` exits with:

```text
Keep requires Accessibility permission.
Open System Settings -> Privacy & Security -> Accessibility
and allow your terminal application or the keep binary.
```

To grant permission:

1. Open `System Settings`.
2. Go to `Privacy & Security`.
3. Open `Accessibility`.
4. Enable your terminal application, for example `Terminal`, `iTerm2`, `Ghostty`, `WezTerm`, or `Alacritty`.
5. Run `keep` again.

If you run the installed binary directly and macOS shows the binary as the requesting process, allow `keep` itself.

## How It Works

Every interval, `keep` asks macOS how long it has been since the last input event.

If the user has been idle for at least the configured interval, `keep`:

1. reads the current cursor position;
2. reads the active display bounds;
3. chooses a nearby valid point on any active display;
4. posts a Core Graphics `kCGEventMouseMoved` event;
5. waits briefly;
6. posts another mouse move event back to the original position.

This handles screen edges and multi-monitor layouts, including displays with negative coordinates.

Because synthetic movement itself resets the macOS idle timer, `keep` tracks its own last synthetic movement and distinguishes it from real user activity.

On Windows, `keep` reads the last input time through `GetLastInputInfo`, enumerates active monitors, and posts a brief absolute mouse move through `SendInput`. After a short delay, it moves the cursor back to its original position. This preserves the same visible mouse-movement approach across platforms while serving the project goal of model loading and downloads, network uptime, and convenience.

## Diagnostics

Run with a short interval while testing:

```bash
keep --interval 5s --verbose
```

Check macOS power assertions:

```bash
pmset -g assertions
```

After synthetic movement, macOS should report recent user activity, commonly through a `UserIsActive` assertion owned by `WindowServer`.

Useful checks:

```bash
keep --version
keep --help
make test
make build
```

## Troubleshooting

### `command not found: keep`

The install path is probably not in `PATH`.

Run directly:

```bash
$HOME/.local/bin/keep
```

Or add it to your shell profile:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Accessibility permission was granted, but movement still fails

Try the following:

- restart the terminal application;
- remove and re-add the terminal application in Accessibility settings;
- if launching `/usr/local/bin/keep` or `$HOME/.local/bin/keep` directly, allow the binary if macOS lists it separately;
- run `keep --interval 5s --verbose` and inspect the logs.

### Cursor moves but a remote desktop session does not see it

Remote desktop clients decide which local input events are forwarded into the remote session. `keep` generates mouse movement in the Windows desktop when launched from Windows or WSL; the remote client must accept and forward that movement.

## Uninstall

If installed into `$HOME/.local/bin`:

```bash
rm -f "$HOME/.local/bin/keep"
```

If installed into `/usr/local/bin`:

```bash
sudo rm -f /usr/local/bin/keep
```

You can also remove the Accessibility permission from `System Settings -> Privacy & Security -> Accessibility`.

## Development

```bash
make test
make build
bin/keep --help
bin/keep --interval 5s --verbose
```

Project layout:

```text
.
├── cmd/keep/main.go
├── internal/activity
├── internal/cli
├── internal/mover
├── go.mod
├── Makefile
├── README.md
└── LICENSE
```

Manual smoke checklist:

- `keep` starts a long-running foreground process.
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

Use `keep` only in ways that comply with your organization policies, contracts, and local rules. This tool is intentionally transparent so its behavior is easy to inspect and reason about.

## License

MIT.
