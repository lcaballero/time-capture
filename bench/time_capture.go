package bench

import (
	"fmt"
	"io"
	"os"
	"time"
)

// A zero parameter function without a return type used to execute arbitrary code
// which can be timed with a TimeCapture.
type Action func()

// A Delta captures a start and stop timestamp.
type Delta struct {
	start     int64
	end       int64
	isStopped bool
	isStarted bool
}

// Captures a start time.  Repeated calls to Start() do not change the captured time.
func (t *Delta) Start() {
	if !t.isStarted {
		t.start = time.Now().UnixNano()
		t.isStarted = true
	}
}

// Captures a stop time.  Repeated calls to Stop() do not change the capture time.
func (t *Delta) Stop() {
	if !t.isStopped {
		t.end = time.Now().UnixNano()
		t.isStopped = true
	}
}

// Provides the difference between start and stop (as a time.Duration).
func (t *Delta) Elapsed() time.Duration {
	t.Stop()
	return time.Duration(t.Diff())
}

// Provides the difference between start and stop as an int64.
func (t *Delta) Diff() int64 {
	return t.end - t.start
}

// Formats the elapsed time in milliseconds.
func (t *Delta) String() string {
	return fmt.Sprintf("Elapsed %d (milliseconds)", t.Elapsed()/time.Millisecond)
}

// Dumps the value for String() to os.Stdout
func (t *Delta) Dump() {
	t.Out(os.Stdout)
}

// Dumps the value for String() to the given writer
func (t *Delta) Out(w io.Writer) {
	fmt.Fprint(w, t.String())
}

// A TimeCapture represents the delta around the execution of a piece of code.
type TimeCapture struct {
	Delta
}

// Start creates and starts a new TimeCapture that it returns.
func Start() *TimeCapture {
	tc := &TimeCapture{}
	tc.Start()
	return tc
}

// Calculates a delta around action and return that TimeCapture.
func Capture(a Action) *TimeCapture {
	tc := Start()
	a()
	tc.Stop()
	return tc
}
