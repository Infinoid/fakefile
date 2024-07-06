package fakefile

import (
	"fmt"
	"io"
)

// TODO: implement io.ByteWriter
// TODO: implement io.StringWriter
// TODO: implement io.WriterTo
type fakefile_writer struct {
	f      *Fakefile
	offset int64
	closed bool
}

// implements io.Writer
func (w *fakefile_writer) Write(b []byte) (int, error) {
	rv, err := w.WriteAt(b, w.offset)
	if err == nil {
		w.offset += int64(rv)
	}
	return rv, err
}

// implements io.Closer
func (w *fakefile_writer) Close() error {
	if w.closed {
		return io.ErrClosedPipe
	}
	w.closed = true
	return nil
}

// implements WriterAt
func (w *fakefile_writer) WriteAt(b []byte, offset int64) (int, error) {
	if w.closed {
		return 0, io.ErrClosedPipe
	}
	w.f.l.Lock()
	defer w.f.l.Unlock()

	if offset+int64(len(b)) > int64(len(w.f.b)) {
		// extend fakefile buffer size
		delta := offset + int64(len(b)) - int64(len(w.f.b))
		zeros := make([]byte, delta)
		w.f.b = append(w.f.b, zeros...)
	}

	rv := copy(w.f.b[offset:], b)

	return rv, nil
}

// implements io.Seeker
// Seeking past the end and then doing a Write will zero-pad up the buffer to the new offset (or run out of memory trying).
func (w *fakefile_writer) Seek(offset int64, whence int) (int64, error) {
	var requested_offset int64
	switch whence {
	case io.SeekStart:
		requested_offset = offset
	case io.SeekCurrent:
		requested_offset = w.offset + offset
	case io.SeekEnd:
		requested_offset = int64(len(w.f.b)) + offset
	default:
		return 0, fmt.Errorf("unknown `whence` value %d", whence)
	}
	if requested_offset < 0 {
		return 0, fmt.Errorf("cannot seek to a negative offset")
	}
	// seeking past the end of the file is not an error; the next write will just pad up the data.
	w.offset = requested_offset
	return requested_offset, nil
}
