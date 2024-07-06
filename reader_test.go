package fakefile

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	fr := fakefile_reader{f: &f}
	r := io.Reader(&fr)
	b := make([]byte, 2)

	rv, err := r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(2), fr.offset)
	require.Equal(t, byte(1), b[0])
	require.Equal(t, byte(2), b[1])

	rv, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 1, rv)
	require.Equal(t, int64(3), fr.offset)
	require.Equal(t, byte(3), b[0])

	_, err = r.Read(b)
	require.Equal(t, io.EOF, err)

	fr.offset = 50
	_, err = r.Read(b)
	require.Equal(t, io.EOF, err)
}

func TestReadAt(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	fr := fakefile_reader{f: &f}
	r := io.ReaderAt(&fr)
	b := make([]byte, 2)

	rv, err := r.ReadAt(b, 0)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, byte(1), b[0])
	require.Equal(t, byte(2), b[1])

	rv, err = r.ReadAt(b, 1)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, byte(2), b[0])
	require.Equal(t, byte(3), b[1])

	rv, err = r.ReadAt(b, 2)
	require.NoError(t, err)
	require.Equal(t, 1, rv)
	require.Equal(t, byte(3), b[0])

	_, err = r.ReadAt(b, 3)
	require.Equal(t, io.EOF, err)

	_, err = r.ReadAt(b, 50)
	require.Equal(t, io.EOF, err)
}

func TestReaderClose(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	r := f.Reader()
	b := make([]byte, 2)

	err := r.Close()
	require.NoError(t, err)
	_, err = r.Read(b)
	require.Error(t, err)
	_, err = r.ReadAt(b, 0)
	require.Error(t, err)
	err = r.Close()
	require.Error(t, err)
}

func TestReaderSeek(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	fr := fakefile_reader{f: &f}
	r := io.ReadSeekCloser(&fr)
	require.NotNil(t, r)

	o, err := r.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), o)

	o, err = r.Seek(1, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(1), o)

	o, err = r.Seek(1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(2), o)

	o, err = r.Seek(1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(3), o)

	o, err = r.Seek(-1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(2), o)

	// Read() uses the new offset
	b := []byte{0}
	i, err := r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 1, i)
	require.Equal(t, byte(3), b[0])
	i, err = r.Read(b)
	require.ErrorIs(t, io.EOF, err)
	require.Equal(t, 0, i)
	require.Equal(t, byte(3), b[0])

	o, err = r.Seek(-2, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(1), o)

	// seeking before the start of the file is not allowed
	_, err = r.Seek(-4, io.SeekEnd)
	require.ErrorContains(t, err, "cannot seek to a negative offset")

	// seeking past the end of the file is allowed
	o, err = r.Seek(2, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(5), o)

	// and a subsequent read just returns EOF
	i, err = r.Read(b)
	require.ErrorIs(t, io.EOF, err)
	require.Equal(t, 0, i)
	require.Equal(t, byte(3), b[0])

	// bad whence value is rejected
	_, err = r.Seek(2, 42)
	require.ErrorContains(t, err, "value 42")
}
