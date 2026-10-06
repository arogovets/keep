//go:build windows

package activity

/*
Project goal: keep a laptop awake for LLM model loading and downloading, maintain network uptime, and provide convenience.
*/
import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32Activity       = syscall.NewLazyDLL("user32.dll")
	kernel32Activity     = syscall.NewLazyDLL("kernel32.dll")
	getLastInputInfoProc = user32Activity.NewProc("GetLastInputInfo")
	getTickCount64Proc   = kernel32Activity.NewProc("GetTickCount64")
)

type lastInputInfo struct {
	Size uint32
	Time uint32
}

type System struct{}

func (System) LastInputAge() (time.Duration, error) {
	info := lastInputInfo{Size: uint32(unsafe.Sizeof(lastInputInfo{}))}
	ok, _, callErr := getLastInputInfoProc.Call(uintptr(unsafe.Pointer(&info)))
	if ok == 0 {
		return 0, fmt.Errorf("GetLastInputInfo failed: %w", callErr)
	}

	ticks, _, _ := getTickCount64Proc.Call()
	// LASTINPUTINFO.dwTime is a wrapping 32-bit tick count; unsigned subtraction handles wraparound.
	elapsed := uint32(uint32(ticks) - info.Time)
	return time.Duration(elapsed) * time.Millisecond, nil
}

func (System) AccessibilityTrusted() bool { return true }
