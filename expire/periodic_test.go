package expire

import (
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeCapture(t *testing.T) {
	Convey("New bell should have a non-null timer", t, func() {
		bell := New(50*time.Millisecond, 200*time.Millisecond, 20*time.Millisecond)
		defer bell.Stop()

		So(bell.C, ShouldNotBeNil)
	})
}
