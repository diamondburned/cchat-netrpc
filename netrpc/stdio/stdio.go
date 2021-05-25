// Package stdio provides a way to gracefully take over both stdin and stdout.
// The package assumes that no other caller may set os.Stdin and os.Stdout,
// which will cause a race condition.
package stdio

import (
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var (
	lock  sync.Mutex
	taken bool

	null    *os.File
	nullErr error
)

func init() {
	null, nullErr = os.OpenFile(os.DevNull, os.O_RDWR|os.O_EXCL, os.ModePerm)
}

// IO describes both the stdin and stdout. It satisfies io.ReadWriteCloser.
type IO struct {
	In  *os.File
	Out *os.File

	closed bool
}

var _ io.ReadWriteCloser = (*IO)(nil)

// ErrStdioTaken is returned if stdio is already takee.
var ErrStdioTaken = errors.New("stdio is taken elsewhere")

// Take takes over the global files to be returned into a ReadWriteCloser. The
// global files will point towards /dev/null until Return is called. If an IO is
// already taken, then it must be closed, otherwise Take will error out.
func Take() (*IO, error) {
	if nullErr != nil {
		return nil, errors.Wrap(nullErr, "failed to create devnull")
	}

	lock.Lock()
	defer lock.Unlock()

	if taken {
		return nil, ErrStdioTaken
	}

	io := IO{
		In:  os.Stdin,
		Out: os.Stdout,
	}

	os.Stdin = null
	os.Stdout = null
	taken = true

	return &io, nil
}

// Read implements io.Reader onto os.Stdin.
func (io *IO) Read(b []byte) (int, error) {
	return io.In.Read(b)
}

// Write implements io.Writer onto os.Stdout.
func (io *IO) Write(b []byte) (int, error) {
	return io.Out.Write(b)
}

// Close returns back the stdio pair to the global scope. The IO called will be
// set to /dev/null. If Close is called multiple times, then the returning only
// happens on the first call.
func (io *IO) Close() error {
	lock.Lock()
	defer lock.Unlock()

	// IO not taken or already closed.
	if !taken || io.closed {
		return os.ErrClosed
	}

	os.Stdin = io.In
	os.Stdout = io.Out

	io.In = null
	io.Out = null
	io.closed = true

	taken = false
	return nil
}
