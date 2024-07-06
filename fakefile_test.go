package fakefile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakefile(t *testing.T) {
	f := New()
	require.Equal(t, 0, len(f.b))

	f = NewLen(0)
	require.Equal(t, 0, len(f.b))

	f = NewLen(3)
	require.Equal(t, 3, len(f.b))

	f = NewFrom([]byte{1, 2})
	require.Equal(t, 2, len(f.b))
	require.Equal(t, byte(1), f.b[0])
	require.Equal(t, byte(2), f.b[1])

	r := f.Reader()
	require.Equal(t, r.f, f)
	require.Equal(t, int64(0), r.offset)

	w := f.Writer()
	require.Equal(t, w.f, f)
	require.Equal(t, int64(0), w.offset)
}
