package mover

import "testing"

func TestTargetPositionInsideDisplay(t *testing.T) {
	got, ok := TargetPosition(Point{X: 10, Y: 10}, 1, []Rect{{X: 0, Y: 0, Width: 100, Height: 100}})
	if !ok {
		t.Fatal("TargetPosition() ok = false, want true")
	}
	if got != (Point{X: 11, Y: 10}) {
		t.Fatalf("TargetPosition() = %+v, want x+1", got)
	}
}

func TestTargetPositionAtRightEdgeFallsBack(t *testing.T) {
	got, ok := TargetPosition(Point{X: 99, Y: 50}, 1, []Rect{{X: 0, Y: 0, Width: 100, Height: 100}})
	if !ok {
		t.Fatal("TargetPosition() ok = false, want true")
	}
	if got != (Point{X: 98, Y: 50}) {
		t.Fatalf("TargetPosition() = %+v, want x-1", got)
	}
}

func TestTargetPositionSupportsNegativeDisplayCoordinates(t *testing.T) {
	displays := []Rect{
		{X: -100, Y: 0, Width: 100, Height: 100},
		{X: 0, Y: 0, Width: 100, Height: 100},
	}
	got, ok := TargetPosition(Point{X: -1, Y: 50}, 1, displays)
	if !ok {
		t.Fatal("TargetPosition() ok = false, want true")
	}
	if got != (Point{X: 0, Y: 50}) {
		t.Fatalf("TargetPosition() = %+v, want cross-display x+1", got)
	}
}

func TestTargetPositionRejectsInvalidDistance(t *testing.T) {
	if _, ok := TargetPosition(Point{X: 1, Y: 1}, 0, []Rect{{X: 0, Y: 0, Width: 10, Height: 10}}); ok {
		t.Fatal("TargetPosition() ok = true, want false")
	}
}
