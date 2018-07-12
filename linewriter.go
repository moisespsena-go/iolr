package ioutil

import (
	"io"
	"os"
)

var (
	StdOutLW = NewLineWriter(os.Stdout)
	StdErrLW = NewLineWriter(os.Stderr)
)

type LineWriter interface {
	io.Writer
	SetCR(v bool)
	CR() bool
	WriteLine() (n int, err error)
	WriteLineB(data []byte) (n int, err error)
	WriteLineS(data string) (n int, err error)
	WriteS(data string) (n int, err error)
}

type DefaultLineWriter struct {
	io.Writer
	cr bool
}

func NewLineWriter(w io.Writer) LineWriter {
	return &DefaultLineWriter{Writer: w}
}

func (d *DefaultLineWriter) SetCR(v bool) {
	d.cr = v
}

func (d *DefaultLineWriter) CR() bool {
	return d.cr
}

func (w *DefaultLineWriter) WriteLine() (n int, err error) {
	if w.cr {
		return w.Write(CRLF)
	} else {
		return w.Write(LF)
	}
}

func (w *DefaultLineWriter) WriteLineB(data []byte) (n int, err error) {
	if w.cr {
		data = append(data, CRLF...)
	} else {
		data = append(data, LF...)
	}
	return w.Write(data)
}

func (w *DefaultLineWriter) WriteLineS(data string) (n int, err error) {
	return w.WriteLineB([]byte(data))
}

func (w *DefaultLineWriter) WriteS(data string) (n int, err error) {
	return w.Write([]byte(data))
}
