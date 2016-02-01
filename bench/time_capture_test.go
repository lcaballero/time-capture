package bench

import (
	"bytes"
	"testing"
	"time"

	"fmt"
	"io/ioutil"
	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeCapture(t *testing.T) {

	Convey("String output should windup in file (instead of os.Stdout) ", t, func() {
		stdout := os.Stdout
		defer func() { os.Stdout = stdout }()

		f, _ := ioutil.TempFile("/tmp", "tc-tests-")
		fmt.Println(f.Name())

		os.Stdout = f
		t := TimeCapture{}
		t.Dump()
		f.Sync()
		os.Stdout = stdout
		f.Close()

		bytes, _ := ioutil.ReadFile(f.Name())
		content := string(bytes)
		str := t.String()

		So(content, ShouldNotEqual, "")
		So(content, ShouldEqual, str)
	})

	Convey("Capture time should out a non-zero string", t, func() {
		delay := 100 * time.Millisecond
		tc := Capture(func() {
			<-time.After(delay)
		})
		buf := bytes.NewBufferString("")
		tc.Out(buf)
		s := buf.String()
		So(s, ShouldNotEqual, "")
	})

	Convey("Capture time of function should take >= delayed time", t, func() {
		delay := 200 * time.Millisecond
		tc := Capture(func() {
			<-time.After(delay)
		})
		So(tc.Elapsed(), ShouldBeGreaterThanOrEqualTo, delay)
	})

	Convey("A started delta should have a start >= to now", t, func() {
		d := &Delta{}
		t := time.Now().UnixNano()
		d.Start()

		So(d.start, ShouldBeGreaterThanOrEqualTo, t)
		So(d.end, ShouldEqual, 0)
		So(d.Elapsed(), ShouldBeGreaterThanOrEqualTo, 0)
	})

	Convey("Should only be able to call Stop() once and capture that timestamp", t, func() {
		d := &Delta{}
		d.Start()
		<-time.NewTimer(50 * time.Millisecond).C
		d.Stop()
		originalStop := d.end
		<-time.NewTimer(50 * time.Millisecond).C
		d.Stop()

		So(d.isStarted, ShouldBeTrue)
		So(d.isStopped, ShouldBeTrue)
		So(d.end, ShouldEqual, originalStop)
	})

	Convey("Should only be able to call Start() once and capture that timestamp", t, func() {
		d := &Delta{}
		d.Start()
		originalStart := d.start
		<-time.NewTimer(50 * time.Millisecond).C
		d.Start()

		So(d.isStopped, ShouldBeFalse)
		So(d.isStarted, ShouldBeTrue)
		So(d.start, ShouldEqual, originalStart)
	})

	Convey("Calling Elapsed() should stop the timer", t, func() {
		d := &Delta{}
		d.Elapsed()
		So(d.isStopped, ShouldBeTrue)
	})

	Convey("New Delta should have ellapsed time of 0", t, func() {
		d := &Delta{}
		So(d.start, ShouldEqual, 0)
		So(d.end, ShouldEqual, 0)
		So(d.isStarted, ShouldBeFalse)
		So(d.isStopped, ShouldBeFalse)
		So(d.Diff(), ShouldEqual, 0)
		So(int64(d.Elapsed()), ShouldBeGreaterThan, 0)
	})
}
