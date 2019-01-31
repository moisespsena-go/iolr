package iolr

import (
	"io"
	"os"
)

var StdinLR = NewLineReader(os.Stdin)

func ReadLine(r io.Reader) (data []byte, err error) {
	b := make([]byte, 1)
	var n int
	for err == nil {
		n, err = r.Read(b)
		if err != nil {
			return
		}
		if n != 1 {
			break
		}
		if b[0] == '\n' {
			break
		}
		data = append(data, b[0])
	}
	return
}

func ReadLineCR(r io.Reader) (data []byte, err error) {
	data, err = ReadLine(r)
	if err != nil {
		return
	}
	if data[len(data)-1] == '\r' {
		data = data[:len(data)-1]
	}
	return
}

type LineReader interface {
	io.Reader
	SetCR(v bool)
	CR() bool
	ReadLine() (data []byte, err error)
	ReadLineS() (data string, err error)
}

type DefaultLineReader struct {
	io.Reader
	cr bool
}

func NewLineReader(r io.Reader) LineReader {
	return &DefaultLineReader{Reader: r}
}

func (r *DefaultLineReader) SetCR(cr bool) {
	r.cr = cr
}

func (r *DefaultLineReader) CR() bool {
	return r.cr
}

func (r *DefaultLineReader) ReadLine() (data []byte, err error) {
	if r.cr {
		return ReadLineCR(r)
	}
	return ReadLine(r)
}

func (r *DefaultLineReader) ReadLineS() (data string, err error) {
	var d []byte
	d, err = r.ReadLine()
	if err == nil {
		data = string(d)
	}
	return
}
