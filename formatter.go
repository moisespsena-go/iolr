package iou

import (
	"strings"

	"fmt"

	"github.com/go-errors/errors"
)

type Formatter interface {
	Write(w LineWriter, require bool, requireMessage *string, defaul string) (err error)
	DefaultValue() interface{}
	Validate(data []byte) error
}

type FormatValuer interface {
	Formatter
	ValueOf(data []byte) interface{}
}

type FOptions struct {
	Message string
	Options []string
	Sep     string
	Default interface{}
	Wrap    [2]string
}

func (f *FOptions) Write(w LineWriter, require bool, requireMessage *string, defaul string) error {
	options := f.Options
	if f.Sep == "" {
		f.Sep = "|"
	}
	if f.Wrap[0] == "" {
		f.Wrap[0] = "("
	}
	if f.Wrap[1] == "" {
		f.Wrap[1] = ")"
	}

	if defaul != "" {
		for i, option := range options {
			if option == defaul {
				options[i] = "*" + option
				break
			}
		}
	}

	if f.Message == "" {
		f.Message = "Choose an option"
	}
	var msg string
	if require {
		msg = *requireMessage
	}
	msg += f.Message + " " + f.Wrap[0] + strings.Join(options, f.Sep) + f.Wrap[1]
	_, err := w.Write([]byte(msg))
	return err
}

func (f *FOptions) Validate(data []byte) error {
	ds := string(data)
	for _, option := range f.Options {
		if option == ds {
			return nil
		}
	}
	return errors.New("Invalid option.")
}

func (f *FOptions) DefaultValue() interface{} {
	return f.Default
}

type FOptionsPairs struct {
	Message        string
	Header         string
	Options        StringPairs
	Default        interface{}
	DefaultWrap    [2]string
	DefaultMessage string
	Sep            string
	CaseSensitive  bool
	optionsMap     map[string]int
}

func (f *FOptionsPairs) init() {
	if f.optionsMap == nil {
		f.optionsMap = map[string]int{}
		for i, pair := range f.Options {
			f.optionsMap[pair.K] = i
		}
	}
}

func (f *FOptionsPairs) Write(w LineWriter, require bool, requireMessage *string, defaul string) (err error) {
	f.init()
	if f.DefaultWrap[0] == "" {
		f.DefaultWrap[0] = "«"
	}
	if f.DefaultWrap[1] == "" {
		f.DefaultWrap[1] = "»"
	}
	if f.DefaultMessage == "" {
		f.DefaultMessage = "(default is %v)"
	}

	if f.Sep == "" {
		f.Sep = ") "
	}

	if f.Header == "" {
		f.Header = "Options:"
	}

	if f.Message == "" {
		f.Message = "Choose an option"
	}

	_, err = w.WriteLineS(f.Header)
	if err != nil {
		return
	}

	var keyv string

	for _, pair := range f.Options {
		if pair.K == "" {
			if _, err = w.WriteLineS(""); err != nil {
				return err
			}
			continue
		}
		keyv = pair.K
		if pair.K == defaul {
			keyv = "*" + pair.K
		}
		_, err = w.WriteLineS(fmt.Sprintf("  %v%v%v", keyv, f.Sep, pair.V))
		if err != nil {
			return
		}
	}

	var msg string
	if require {
		msg = *requireMessage
	}
	msg += f.Message

	if defaul != "" {
		msg += " " + fmt.Sprintf(f.DefaultMessage, f.DefaultWrap[0]+defaul+f.DefaultWrap[1])
	}
	_, err = w.WriteS(msg)
	return
}

func (f *FOptionsPairs) Validate(data []byte) error {
	f.init()
	ds := string(data)

	if f.CaseSensitive {
		if _, ok := f.optionsMap[ds]; ok {
			return nil
		}
	} else {
		ds = strings.ToLower(ds)
		for _, pair := range f.Options {
			if strings.ToLower(pair.K) == ds {
				return nil
			}
		}
	}
	return errors.New("Invalid option.")
}

func (f *FOptionsPairs) DefaultValue() interface{} {
	return f.Default
}

func (f *FOptionsPairs) ValueOf(data []byte) (value interface{}) {
	f.init()
	if i, ok := f.optionsMap[string(data)]; ok {
		value = f.Options[i].Ki
	}
	return
}
