package finch

import (
	"math"
	"time"
)

const MaxFixedFrames = 5

type Time struct {
	targetMS    float64
	startMS     float64
	currentMS   float64
	deltaMS     float64
	elapsedMS   float64
	fixedFrames int
}

func NewTime(targetFPS float64) *Time {
	return &Time{
		targetMS: 1000.0 / targetFPS,
	}
}

func (t *Time) start() {
	now := float64(time.Now().UnixMilli())

	t.startMS = now
	t.currentMS = now
}

func (t *Time) tick() {
	now := float64(time.Now().UnixMilli())
	prev := t.currentMS

	t.currentMS = now
	t.deltaMS = now - prev

	if t.deltaMS < 0 {
		t.deltaMS = 0
	}

	t.elapsedMS += t.deltaMS

	t.fixedFrames = int(math.Floor(t.elapsedMS / t.targetMS))
	if t.fixedFrames > 0 {
		if t.fixedFrames > MaxFixedFrames {
			t.fixedFrames = MaxFixedFrames
		}
		t.elapsedMS -= float64(t.fixedFrames) * t.targetMS
	}
}

func (t *Time) DeltaMilli() float64 {
	return t.deltaMS
}

func (t *Time) DeltaSeconds() float64 {
	return t.deltaMS / 1000.0
}

func (t *Time) FixedMilli() float64 {
	return t.targetMS
}

func (t *Time) FixedSeconds() float64 {
	return t.targetMS / 1000.0
}

func (t *Time) FixedFrames() int {
	return t.fixedFrames
}
