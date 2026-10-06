//go:build windows

package mover

/*
Project goal: keep a laptop awake for LLM model loading and downloading, maintain network uptime, and provide convenience.
Windows uses the same small, temporary mouse movement as macOS so the system sees input activity.
*/
import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	inputMouse            = 0
	mouseEventMove        = 0x0001
	mouseEventAbsolute    = 0x8000
	mouseEventVirtualDesk = 0x4000
	smVirtualScreenLeft   = 76
	smVirtualScreenTop    = 77
	smVirtualScreenWidth  = 78
	smVirtualScreenHeight = 79
	returnDelay           = 100 * time.Millisecond
)

type point struct{ X, Y int32 }
type nativeRect struct{ Left, Top, Right, Bottom int32 }
type monitorInfo struct {
	Size    uint32
	Monitor nativeRect
	Work    nativeRect
	Flags   uint32
}
type mouseInput struct {
	Dx, Dy    int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}
type nativeInput struct {
	Type  uint32
	Mouse mouseInput
}

var (
	user32Mouse          = syscall.NewLazyDLL("user32.dll")
	getCursorPosProc     = user32Mouse.NewProc("GetCursorPos")
	sendInputProc        = user32Mouse.NewProc("SendInput")
	getSystemMetricsProc = user32Mouse.NewProc("GetSystemMetrics")
	enumDisplayMonitors  = user32Mouse.NewProc("EnumDisplayMonitors")
	getMonitorInfoProc   = user32Mouse.NewProc("GetMonitorInfoW")
)

type System struct{}

func (System) StartupDetail(distance int) string {
	return fmt.Sprintf("Move distance: %d screen pixel(s)", distance)
}

func (System) ActionDescription() string { return "cursor moved" }

func (System) MoveAndReturn(ctx context.Context, distance int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	position, err := CurrentPosition()
	if err != nil {
		return err
	}
	displays, err := ActiveDisplays()
	if err != nil {
		return err
	}
	target, ok := TargetPosition(position, distance, displays)
	if !ok {
		return errors.New("no valid target position within active displays")
	}
	if err := postMove(target); err != nil {
		return err
	}

	timer := time.NewTimer(returnDelay)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
	if err := postMove(position); err != nil {
		return err
	}
	return nil
}

func CurrentPosition() (Point, error) {
	var p point
	ok, _, callErr := getCursorPosProc.Call(uintptr(unsafe.Pointer(&p)))
	if ok == 0 {
		return Point{}, fmt.Errorf("GetCursorPos failed: %w", callErr)
	}
	return Point{X: float64(p.X), Y: float64(p.Y)}, nil
}

func ActiveDisplays() ([]Rect, error) {
	var rects []Rect
	callback := syscall.NewCallback(func(monitor uintptr, _, _ uintptr, _ uintptr) uintptr {
		info := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
		ok, _, _ := getMonitorInfoProc.Call(monitor, uintptr(unsafe.Pointer(&info)))
		if ok != 0 {
			rects = append(rects, Rect{
				X: float64(info.Monitor.Left), Y: float64(info.Monitor.Top),
				Width:  float64(info.Monitor.Right - info.Monitor.Left),
				Height: float64(info.Monitor.Bottom - info.Monitor.Top),
			})
		}
		return 1
	})
	ok, _, callErr := enumDisplayMonitors.Call(0, 0, callback, 0)
	if ok == 0 {
		return nil, fmt.Errorf("EnumDisplayMonitors failed: %w", callErr)
	}
	if len(rects) == 0 {
		return nil, errors.New("no active displays found")
	}
	return rects, nil
}

func postMove(p Point) error {
	left := int32(getMetric(smVirtualScreenLeft))
	top := int32(getMetric(smVirtualScreenTop))
	width := int32(getMetric(smVirtualScreenWidth))
	height := int32(getMetric(smVirtualScreenHeight))
	if width <= 1 || height <= 1 {
		return errors.New("invalid virtual screen bounds")
	}
	x := int32((p.X-float64(left))*65535/float64(width-1) + 0.5)
	y := int32((p.Y-float64(top))*65535/float64(height-1) + 0.5)
	input := nativeInput{Type: inputMouse, Mouse: mouseInput{
		Dx: x, Dy: y, Flags: mouseEventMove | mouseEventAbsolute | mouseEventVirtualDesk,
	}}
	count, _, callErr := sendInputProc.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	if count != 1 {
		return fmt.Errorf("SendInput failed to move cursor: %w", callErr)
	}
	return nil
}

func getMetric(index int32) int {
	value, _, _ := getSystemMetricsProc.Call(uintptr(index))
	return int(int32(value))
}

func (p Point) String() string { return fmt.Sprintf("(%.0f, %.0f)", p.X, p.Y) }
