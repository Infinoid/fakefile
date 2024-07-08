package fakefile

import "sync"

type Fakefile struct {
	b []byte
	l sync.RWMutex
}

func New() *Fakefile {
	return NewLen(0)
}

func NewLen(n int) *Fakefile {
	return &Fakefile{b: make([]byte, n)}
}

func NewFrom(b []byte) *Fakefile {
	return &Fakefile{b: b}
}

func (f *Fakefile) Reader() *fakefile_reader {
	return &fakefile_reader{f: f}
}

func (f *Fakefile) Writer() *fakefile_writer {
	return &fakefile_writer{f: f}
}

func (f *Fakefile) Bytes() []byte {
	return f.b
}
