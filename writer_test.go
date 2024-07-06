package fakefile

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{0, 0, 0}}
	fr := fakefile_writer{f: &f}
	r := io.Writer(&fr)
	b := []byte{1, 2}

	rv, err := r.Write(b)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(2), fr.offset)
	require.Equal(t, 3, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(2), f.b[1])
	require.Equal(t, byte(0), f.b[2])

	rv, err = r.Write(b)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(4), fr.offset)
	require.Equal(t, 4, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(2), f.b[1])
	require.Equal(t, byte(1), f.b[2])
	require.Equal(t, byte(2), f.b[3])
}

func TestWriteAt(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{0, 0, 0}}
	fr := fakefile_writer{f: &f}
	r := io.WriterAt(&fr)
	b := []byte{1, 2}

	rv, err := r.WriteAt(b, 0)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(0), fr.offset)
	require.Equal(t, 3, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(2), f.b[1])
	require.Equal(t, byte(0), f.b[2])

	rv, err = r.WriteAt(b, 1)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(0), fr.offset)
	require.Equal(t, 3, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(1), f.b[1])
	require.Equal(t, byte(2), f.b[2])

	rv, err = r.WriteAt(b, 5)
	require.NoError(t, err)
	require.Equal(t, 2, rv)
	require.Equal(t, int64(0), fr.offset)
	require.Equal(t, 7, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(1), f.b[1])
	require.Equal(t, byte(2), f.b[2])
	require.Equal(t, byte(0), f.b[3])
	require.Equal(t, byte(0), f.b[4])
	require.Equal(t, byte(1), f.b[5])
	require.Equal(t, byte(2), f.b[6])
}

func TestWriterClose(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	r := fakefile_writer{f: &f}
	b := make([]byte, 2)

	err := r.Close()
	require.NoError(t, err)
	_, err = r.Write(b)
	require.Error(t, err)
	_, err = r.WriteAt(b, 0)
	require.Error(t, err)
	err = r.Close()
	require.Error(t, err)
}

func TestWriterSeek(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	f := Fakefile{b: []byte{1, 2, 3}}
	fw := fakefile_writer{f: &f}
	w := io.WriteSeeker(&fw)
	require.NotNil(t, w)

	o, err := w.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), o)

	o, err = w.Seek(1, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(1), o)

	o, err = w.Seek(1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(2), o)

	o, err = w.Seek(1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(3), o)

	o, err = w.Seek(-1, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(2), o)

	// Write() uses the new offset
	b := []byte{4}
	i, err := w.Write(b)
	require.NoError(t, err)
	require.Equal(t, 1, i)
	require.Equal(t, byte(4), f.b[2])
	// Write() at the end of the file extends the file
	i, err = w.Write(b)
	require.NoError(t, err)
	require.Equal(t, 1, i)
	require.Equal(t, 4, len(f.b))
	require.Equal(t, byte(4), f.b[3])

	o, err = w.Seek(-2, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(2), o)
	require.Equal(t, o, fw.offset)

	// seeking before the start of the file is not allowed
	_, err = w.Seek(-8, io.SeekEnd)
	require.ErrorContains(t, err, "cannot seek to a negative offset")

	// seeking past the end of the file is allowed
	o, err = w.Seek(2, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(6), o)

	// and a subsequent write extends the file size
	require.Equal(t, 4, len(f.b))
	i, err = w.Write(b)
	require.NoError(t, err)
	require.Equal(t, 1, i)
	require.Equal(t, 7, len(f.b))
	require.Equal(t, byte(4), f.b[6])

	// bad whence value is rejected
	_, err = w.Seek(2, 42)
	require.ErrorContains(t, err, "value 42")
}
