# stay

`stay` is a small macOS CLI utility that prevents user idle by posting real mouse movement events through Core Graphics. It runs as a transparent foreground terminal process until you stop it with `Ctrl+C`.

No `.app`, menu bar item, tray icon, daemon, launch agent, AppleScript, or GUI dependency is created.

The implementation is inspired by the general idea of periodic mouse movement tools, but no source code from `automatic-mouse-mover` is copied into this repository.

## Requirements

- macOS on Apple Silicon or Intel.
- Accessibility permission for the terminal application running `stay`, or for the installed `stay` binary.
- Go 1.22 or newer for local builds.

## Install

```bash
make build
make install
```

By default `make install` copies the binary to `$HOME/.local/bin/stay`. Make sure that directory is in your `PATH`.

To install into `/usr/local/bin` instead:

```bash
sudo make install PREFIX=/usr/local
```

## Run

```bash
stay
```

Example output:

```text
13:42:00 Stay started
Interval: 60s
Move distance: 1px
Press Ctrl+C to stop

13:43:00 idle for 60s - cursor moved
13:44:00 user activity detected - skipped
```

While the process is running in the terminal, it is active. Press `Ctrl+C` to stop it. After shutdown it no longer generates events.

## Options

```bash
stay --interval 60s
stay --distance 1
stay --always
stay --verbose
stay --version
stay --help
```

- `--interval` sets the check period. Values accepted by Go duration syntax work, for example `60s`, `1m`, or `500ms`. A bare integer is treated as seconds.
- `--distance` sets the cursor movement distance in pixels.
- `--always` moves the cursor every interval regardless of user activity.
- `--verbose` prints more detailed diagnostic logs.
- `--version` prints the binary version.
- `--help` prints usage.

## Accessibility Permission

If permission is missing, `stay` exits with:

```text
Stay requires Accessibility permission.
Open System Settings -> Privacy & Security -> Accessibility
and allow your terminal application or the stay binary.
```

Grant permission in macOS System Settings, then run `stay` again. No GUI interaction is required after the permission is granted.

## Diagnostics

Useful checks:

```bash
stay --interval 5s --verbose
pmset -g assertions
```

The utility logs startup, detected user activity, each cursor movement, errors, and shutdown. Each synthetic movement moves from the current cursor position to a nearby valid point, waits briefly, and moves back to the original position to avoid accumulated cursor drift.

## Citrix

Synthetic mouse movement reaches Citrix Workspace only under the same conditions in which Citrix accepts normal local cursor movement. In practice this can depend on Citrix focus, session mode, client settings, and the active remote window. Verify whether the Citrix window must be focused for your environment.

## Uninstall

Remove the installed binary:

```bash
rm -f "$HOME/.local/bin/stay"
```

If installed with `/usr/local`:

```bash
sudo rm -f /usr/local/bin/stay
```

You can also remove the Accessibility permission from System Settings.

## Policy

Use this utility only in ways that comply with your organization policies and local rules.

## Development

```bash
make test
make build
bin/stay --help
```

Manual smoke checklist:

- `Ctrl+C` stops the foreground process.
- Cursor returns to its original position after movement.
- No accumulated cursor drift after repeated intervals.
- Movement works on the primary monitor and an additional monitor.
- Cursor near screen edges still moves to a valid on-screen point and returns.
- Missing Accessibility permission produces the documented error.
- `pmset -g assertions` or another system indicator reflects user activity after synthetic movement.
- Citrix behavior is verified with and without active Citrix focus.
