package ioutil

import (
	"io"

	"github.com/moisespsena/go-error-wrap"
)

var (
	STDMessageLR       = NewMessageLineReader(StdOutLW, StdinLR, StdErrLW)
	STDStringMessageLR = NewStringMessageLineReader(StdOutLW, StdinLR, StdErrLW)
)

type MessageLineReader struct {
	Writer         LineWriter
	Reader         LineReader
	ErrorWriter    LineWriter
	Sep            string
	RequireMessage string
}

func NewMessageLineReader(w io.Writer, r io.Reader, ew io.Writer) *MessageLineReader {
	var (
		lw, lew LineWriter
		lr      LineReader
		ok      bool
	)
	if ew == nil {
		ew = w
	}
	if lw, ok = w.(LineWriter); !ok {
		lw = NewLineWriter(w)
	}
	if lew, ok = ew.(LineWriter); !ok {
		lew = NewLineWriter(ew)
	}
	if lr, ok = r.(LineReader); !ok {
		lr = NewLineReader(r)
	}
	return &MessageLineReader{lw, lr, lew, ": ", "* "}
}

func (r *MessageLineReader) ReadRaw(message []byte) (data []byte, err errwrap.ErrorWrapperInterface) {
	_, err2 := r.Writer.Write(message)
	if err2 != nil {
		return nil, errwrap.Wrap(err, "Write Message")
	}
	data, err2 = r.Reader.ReadLine()
	if err2 != nil {
		return data, errwrap.Wrap(err, "Read Line")
	}
	return
}

func (r *MessageLineReader) Read(message string, defaul ...string) (data []byte, err errwrap.ErrorWrapperInterface) {
	data, err = r.ReadRaw(append([]byte(message), []byte(r.Sep)...))
	if err == nil && len(data) == 0 && (len(defaul) > 0 && defaul[0] != "") {
		return []byte(defaul[0]), nil
	}
	return
}

func (r *MessageLineReader) ReadFormatter(formatter Formatter, require bool, defaul ...string) (data []byte, err errwrap.ErrorWrapperInterface) {
	if len(defaul) == 0 || defaul[0] == "" {
		defaul = []string{formatter.DefaultValue()}
	}
	var err2 error
	for {
		err2 = formatter.Write(r.Writer, require, &r.RequireMessage, defaul[0])
		if err2 != nil {
			return nil, errwrap.Wrap(err2, "Write Message")
		}
		data, err = r.ReadRaw([]byte(r.Sep))
		if err != nil {
			break
		}
		if len(data) == 0 {
			if defaul[0] != "" {
				return []byte(defaul[0]), nil
			}
		} else {
			err2 = formatter.Validate(data)
			if err2 != nil {
				_, err2 = r.ErrorWriter.Write([]byte(err2.Error()))
				if err2 != nil {
					return nil, errwrap.Wrap(err2, "Write Error")
				}
				_, err2 = r.ErrorWriter.WriteLine()
				if err2 != nil {
					return nil, errwrap.Wrap(err2, "Write Error")
				}
				continue
			}
		}
		break
	}
	return
}

func (r *MessageLineReader) ReadF(formatter Formatter, defaul ...string) (data []byte, err errwrap.ErrorWrapperInterface) {
	return r.ReadFormatter(formatter, false, defaul...)
}

func (r *MessageLineReader) Require(message string, defaul ...string) (data []byte, err errwrap.ErrorWrapperInterface) {
	for len(data) == 0 && err == nil {
		data, err = r.Read(message, defaul...)
	}
	return
}

func (r *MessageLineReader) RequireF(formatter Formatter, defaul ...string) (data []byte, err errwrap.ErrorWrapperInterface) {
	for len(data) == 0 && err == nil {
		data, err = r.ReadFormatter(formatter, true, defaul...)
	}
	return
}

type StringMessageLineReader struct {
	Reader *MessageLineReader
}

func NewStringMessageLineReader(w io.Writer, r io.Reader, ew io.Writer) *StringMessageLineReader {
	return &StringMessageLineReader{NewMessageLineReader(w, r, ew)}
}

func (r *StringMessageLineReader) ReadRaw(message []byte) (data string, err errwrap.ErrorWrapperInterface) {
	var d []byte
	d, err = r.Reader.ReadRaw(message)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *StringMessageLineReader) Read(message string, defaul ...string) (data string, err errwrap.ErrorWrapperInterface) {
	var d []byte
	d, err = r.Reader.Read(message, defaul...)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *StringMessageLineReader) ReadF(formatter Formatter, defaul ...string) (data string, err errwrap.ErrorWrapperInterface) {
	var d []byte
	d, err = r.Reader.ReadF(formatter, defaul...)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *StringMessageLineReader) Require(message string, defaul ...string) (data string, err errwrap.ErrorWrapperInterface) {
	var d []byte
	d, err = r.Reader.Require(message, defaul...)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *StringMessageLineReader) RequireF(formatter Formatter, defaul ...string) (data string, err errwrap.ErrorWrapperInterface) {
	var d []byte
	d, err = r.Reader.RequireF(formatter, defaul...)
	if err == nil {
		data = string(d)
	}
	return
}
