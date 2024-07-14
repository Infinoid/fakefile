# Fakefile

![Build](https://github.com/Infinoid/fakefile/actions/workflows/go.yaml/badge.svg)
<a href="https://pkg.go.dev/github.com/infinoid/fakefile"><img src="https://pkg.go.dev/badge/github.com/infinoid/fakefile.svg" alt="Docs"></a>

This is a simple in-memory fake file, for testing.  It is like `bytes.Buffer`
but with a few more features:

* concurrent readers/writers
* implements the `io.ReaderAt` interface
* implements the `io.WriterAt` interface
* implements the `io.Seeker` interface

# Usage

At the top:

```go
import "github.com/infinoid/fakefile"
```

Then, in a function:
```go
	ff := fakefile.New()
	w := ff.Writer()
	defer w.Close()
	_, _ = w.Write([]byte("Hello world!\n"))
	_, _ = w.WriteAt([]byte("WORLD"), 6)

	r := ff.Reader()
	defer r.Close()
	all, _ := io.ReadAll(r)
	print(string(all)) // prints "Hello WORLD!\n"
```
