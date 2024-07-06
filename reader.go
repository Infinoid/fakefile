package fakefile

import (
	"fmt"
	"io"
)

// TODO: implement io.ByteReader
// TODO: implement io.ByteScanner
// TODO: implement io.ReaderFrom
// TODO: implement io.RuneReader
// TODO: implement io.RuneScanner
type fakefile_reader struct {
	f      *Fakefile
	offset int64
	closed bool
}

// implements io.Reader
func (r *fakefile_reader) Read(b []byte) (int, error) {
	rv, err := r.ReadAt(b, r.offset)
	if err == nil {
		r.offset += int64(rv)
	}
	return rv, err
}

// implements io.Closer
func (r *fakefile_reader) Close() error {
	if r.closed {
		return io.ErrClosedPipe
	}
	r.closed = true
	return nil
}

// implements io.ReaderAt
func (r *fakefile_reader) ReadAt(b []byte, offset int64) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	r.f.l.RLock()
	defer r.f.l.RUnlock()

	if offset >= int64(len(r.f.b)) {
		return 0, io.EOF
	}

	rv := copy(b, r.f.b[offset:])

	return rv, nil
}

// implements io.Seeker
// Seeking past the end and then doing a Read will return io.EOF.
func (r *fakefile_reader) Seek(offset int64, whence int) (int64, error) {
	var requested_offset int64
	switch whence {
	case io.SeekStart:
		requested_offset = offset
	case io.SeekCurrent:
		requested_offset = r.offset + offset
	case io.SeekEnd:
		requested_offset = int64(len(r.f.b)) + offset
	default:
		return 0, fmt.Errorf("unknown `whence` value %d", whence)
	}
	if requested_offset < 0 {
		return 0, fmt.Errorf("cannot seek to a negative offset")
	}
	// seeking past the end of the file is not an error; the next read will just return EOF.
	r.offset = requested_offset
	return requested_offset, nil
}
