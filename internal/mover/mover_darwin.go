//go:build darwin

package mover

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
#include <stdbool.h>
#include <stdlib.h>

typedef struct {
	double x;
	double y;
	double width;
	double height;
} display_rect;

typedef struct {
	double x;
	double y;
	bool ok;
} mouse_location;

static mouse_location current_mouse_location(void) {
	CGEventRef event = CGEventCreate(NULL);
	if (event == NULL) {
		mouse_location failed = {0, 0, false};
		return failed;
	}
	CGPoint point = CGEventGetLocation(event);
	CFRelease(event);
	mouse_location location = {point.x, point.y, true};
	return location;
}

static bool post_mouse_move(double x, double y) {
	CGPoint point = CGPointMake(x, y);
	CGEventRef event = CGEventCreateMouseEvent(NULL, kCGEventMouseMoved, point, kCGMouseButtonLeft);
	if (event == NULL) {
		return false;
	}
	CGEventPost(kCGHIDEventTap, event);
	CFRelease(event);
	return true;
}

static int active_display_count(void) {
	uint32_t count = 0;
	CGGetActiveDisplayList(0, NULL, &count);
	return (int)count;
}

static int active_display_bounds(display_rect *rects, int max) {
	uint32_t count = 0;
	CGGetActiveDisplayList(0, NULL, &count);
	if (count == 0 || max <= 0) {
		return 0;
	}

	CGDirectDisplayID *displays = malloc(sizeof(CGDirectDisplayID) * count);
	if (displays == NULL) {
		return 0;
	}
	if (CGGetActiveDisplayList(count, displays, &count) != kCGErrorSuccess) {
		free(displays);
		return 0;
	}

	int limit = (int)count;
	if (limit > max) {
		limit = max;
	}
	for (int i = 0; i < limit; i++) {
		CGRect bounds = CGDisplayBounds(displays[i]);
		rects[i].x = bounds.origin.x;
		rects[i].y = bounds.origin.y;
		rects[i].width = bounds.size.width;
		rects[i].height = bounds.size.height;
	}
	free(displays);
	return limit;
}
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unsafe"
)

const returnDelay = 100 * time.Millisecond

type System struct{}

func (System) StartupDetail(distance int) string {
	return fmt.Sprintf("Move distance: %d display point(s)", distance)
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
	point := C.current_mouse_location()
	if !bool(point.ok) {
		return Point{}, errors.New("failed to read current mouse location")
	}
	return Point{X: float64(point.x), Y: float64(point.y)}, nil
}

func ActiveDisplays() ([]Rect, error) {
	count := int(C.active_display_count())
	if count <= 0 {
		return nil, errors.New("no active displays found")
	}

	rects := make([]C.display_rect, count)
	read := int(C.active_display_bounds((*C.display_rect)(unsafe.Pointer(&rects[0])), C.int(count)))
	if read <= 0 {
		return nil, errors.New("failed to read active display bounds")
	}

	displays := make([]Rect, read)
	for i := 0; i < read; i++ {
		displays[i] = Rect{
			X:      float64(rects[i].x),
			Y:      float64(rects[i].y),
			Width:  float64(rects[i].width),
			Height: float64(rects[i].height),
		}
	}
	return displays, nil
}

func postMove(point Point) error {
	if !bool(C.post_mouse_move(C.double(point.X), C.double(point.Y))) {
		return fmt.Errorf("failed to create mouse move event for %s", point)
	}
	return nil
}

func (p Point) String() string {
	return fmt.Sprintf("(%.0f, %.0f)", p.X, p.Y)
}
