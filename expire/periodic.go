package expire

import "time"

// An Periodic is a combination of a ticker and utility functions
// determining when to do extra work, and when to stream/flush
// accumulated work.
type Periodic struct {
	C       chan time.Time
	c       <-chan time.Time
	done    chan bool
	ticker  *time.Ticker
	last    int64
	overdue time.Duration
	fuzz    time.Duration
	time    time.Time
	stopped bool
}

// New creates a Periodic with a ticker set to the recheck duration
// and with the overdue duration provided.  Fuzz is used reduce latency
// by providing a window (+/- fuzz) around the expiration time.
// All parameters are in milliseconds.
func New(recheckMs, overdueMs, fuzzMs time.Duration) *Periodic {
	recheck, overdue, fuzz :=
		recheckMs*time.Millisecond,
		overdueMs*time.Second,
		fuzzMs*time.Millisecond
	tic := time.NewTicker(recheck)
	bell := &Periodic{
		C:       make(chan time.Time, 1),
		c:       tic.C,
		last:    time.Now().UnixNano(),
		ticker:  tic,
		overdue: overdue,
		fuzz:    fuzz,
		done:    make(chan bool, 1),
		time:    time.Now(),
		stopped: false,
	}
	go func() {
		for {
			select {
			case bell.time = <-bell.c:
				bell.C <- bell.time
			case <-bell.done:
				return
			}
		}
	}()
	return bell
}

// Stop turn off the Periodic.  After Stop no more ticks will be
// sent via the ticker. Stop does not close the channel to prevent
// a read from the channel to succeed incorrectly.
func (e *Periodic) Stop() {
	if e.stopped {
		return
	}
	e.stopped = true
	e.ticker.Stop()
	e.done <- true
}

// Determines if more time the the overdue duration has elapsed.
func (e *Periodic) IsOverdue() bool {
	now := e.time.UnixNano()
	delta := e.last - now

	return time.Duration(delta) < (e.overdue + e.fuzz)
}

// Reset starts the expiration period over.
func (e *Periodic) Reset() {
	e.last = time.Now().UnixNano()
}
