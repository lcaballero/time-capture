package bench



type Recorder func(tc *TimeCapture)
type Timed func(a Action)
type Counter func(a int64)

func Time(stat Recorder) Timed {
	return func(a Action) {
		stat(Capture(a))
	}
}

