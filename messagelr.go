package iolr

import (
	"fmt"
	"io"
	"strings"

	"github.com/moisespsena/go-error-wrap"
)

var (
	STDMessageLR = NewMessageLineReader(StdOutLW, StdinLR, StdErrLW, StdErrLW)
)

type MessageLineReader struct {
	Writer             LineWriter
	Reader             LineReader
	ErrorWriter        LineWriter
	InputMessageWriter LineWriter
	Sep                string
	RequireMessage     string
	printInput         bool
	InputMessage       func(data []byte) []byte
}

func NewMessageLineReader(w io.Writer, r io.Reader, ew io.Writer, iw ...io.Writer) *MessageLineReader {
	var (
		lw, lew, liw LineWriter
		lr           LineReader
		ok           bool
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
	if len(iw) > 0 && iw[0] != nil {
		liw = NewLineWriter(iw[0])
	}
	return &MessageLineReader{lw, lr, lew, liw, ": ", "* ", false, func(data []byte) []byte {
		var sufix = "\n"
		if lr.CR() {
			sufix = "\r\n"
		}
		return []byte("« " + strings.TrimSuffix(string(data), sufix) + " »\n\n")
	}}
}

func (r *MessageLineReader) EnablePrintInput() *MessageLineReader {
	r.printInput = true
	return r
}

func (r *MessageLineReader) DisablePrintInput() *MessageLineReader {
	r.printInput = false
	return r
}

func (r *MessageLineReader) IsPrintInputEnabled() bool {
	return r.printInput
}

func (r *MessageLineReader) WithPrintInput() func() {
	if r.printInput {
		return func() {}
	}
	r.printInput = true
	return func() {
		r.printInput = false
	}
}

func (r *MessageLineReader) ReadRaw(message []byte) (data []byte, err errwrap.ErrorWrapper) {
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

func (r *MessageLineReader) read(message string, defaul ...string) (data []byte, err errwrap.ErrorWrapper) {
	data, err = r.ReadRaw(append([]byte(message), []byte(r.Sep)...))
	if err == nil && len(data) == 0 && (len(defaul) > 0 && defaul[0] != "") {
		return []byte(defaul[0]), nil
	}
	return
}

func (r *MessageLineReader) Read(message string, defaul ...string) (data []byte, err errwrap.ErrorWrapper) {
	if data, err = r.read(message, defaul...); err == nil && r.printInput && r.InputMessageWriter != nil {
		r.InputMessageWriter.Write(r.InputMessage(data))
	}
	return
}

func (r *MessageLineReader) ReadFormatter(formatter Formatter, require bool, defaul ...interface{}) (value interface{}, err errwrap.ErrorWrapper) {
	if len(defaul) == 0 || defaul[0] == nil {
		defaul = []interface{}{formatter.DefaultValue()}
	}

	var (
		defaultString string
		err2          error
		data          []byte
	)

	if defaul[0] != nil {
		defaultString = fmt.Sprint(defaul[0])
	}

	for {
		err2 = formatter.Write(r.Writer, require, &r.RequireMessage, defaultString)
		if err2 != nil {
			return nil, errwrap.Wrap(err2, "Write Message")
		}
		data, err = r.ReadRaw([]byte(r.Sep))
		if err != nil {
			break
		}
		if len(data) == 0 {
			for _, value = range defaul {
				if value != nil {
					return
				}
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

			if r.printInput && r.InputMessageWriter != nil {
				r.InputMessageWriter.WriteLineB(r.InputMessage(data))
			}

			if valuer, ok := formatter.(FormatValuer); ok {
				value = valuer.ValueOf(data)
			} else {
				value = data
			}
		}
		break
	}
	return
}

func (r *MessageLineReader) ReadF(formatter Formatter, defaul ...interface{}) (value interface{}, err errwrap.ErrorWrapper) {
	return r.ReadFormatter(formatter, false, defaul...)
}

func (r *MessageLineReader) Require(message string, defaul ...string) (data []byte, err errwrap.ErrorWrapper) {
	for len(data) == 0 && err == nil {
		data, err = r.Read(message, defaul...)
	}
	return
}

func (r *MessageLineReader) RequireF(formatter Formatter, defaul ...interface{}) (value interface{}, err errwrap.ErrorWrapper) {
	for value == nil && err == nil {
		value, err = r.ReadFormatter(formatter, true, defaul...)
	}
	return
}

func (r *MessageLineReader) ReadRawS(message []byte) (data string, err errwrap.ErrorWrapper) {
	var d []byte
	d, err = r.ReadRaw(message)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *MessageLineReader) readS(reader func(message string, defaul ...string) (value []byte, err errwrap.ErrorWrapper),
	message string, defaul ...string) (data string, err errwrap.ErrorWrapper) {
	var d []byte
	d, err = reader(message, defaul...)
	if err == nil {
		data = string(d)
	}
	return
}

func (r *MessageLineReader) ReadS(message string, defaul ...string) (data string, err errwrap.ErrorWrapper) {
	return r.readS(r.Read, message, defaul...)
}

func (r *MessageLineReader) RequireS(message string, defaul ...string) (data string, err errwrap.ErrorWrapper) {
	return r.readS(r.Require, message, defaul...)
}

func (r *MessageLineReader) readStringF(reader func(formatter Formatter, defaul ...interface{}) (value interface{}, err errwrap.ErrorWrapper),
	formatter Formatter, defaul ...string) (value string, err errwrap.ErrorWrapper) {
	var (
		d                interface{}
		defaultInterface = make([]interface{}, len(defaul))
	)

	for i, v := range defaul {
		defaultInterface[i] = v
	}
	d, err = reader(formatter, defaultInterface...)
	if err == nil {
		if d == nil {
			value = ""
			return
		}
		switch dt := d.(type) {
		case []byte:
			value = string(dt)
		case string:
			value = dt
		default:
			value = fmt.Sprint(d)
		}
	}
	return
}

func (r *MessageLineReader) ReadFS(formatter Formatter, defaul ...string) (value string, err errwrap.ErrorWrapper) {
	return r.readStringF(r.ReadF, formatter, defaul...)
}

func (r *MessageLineReader) RequireFS(formatter Formatter, defaul ...string) (data string, err errwrap.ErrorWrapper) {
	return r.readStringF(r.RequireF, formatter, defaul...)
}

func IsEmptyInput(value interface{}) bool {
	return false
}
