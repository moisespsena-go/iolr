package ioutil

import (
	"strings"

	"fmt"
	"sort"

	"github.com/go-errors/errors"
)

type Formatter interface {
	Write(w LineWriter, require bool, requireMessage *string, defaul string) (err error)
	DefaultValue() string
	Validate(data []byte) error
}

type FOptions struct {
	Message string
	Options []string
	Sep     string
	Default string
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

func (f *FOptions) DefaultValue() string {
	return f.Default
}

type FOptionsMap struct {
	Message        string
	Header         string
	Options        map[string]string
	Default        string
	DefaultWrap    [2]string
	DefaultMessage string
	Sep            string
	CaseSensitive  bool
}

func (f *FOptionsMap) Write(w LineWriter, require bool, requireMessage *string, defaul string) (err error) {
	var keys []string
	for key := range f.Options {
		keys = append(keys, key)
	}
	sort.Strings(keys)

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

	for _, key := range keys {
		keyv = key
		if key == defaul {
			keyv = "*" + key
		}
		_, err = w.WriteLineS(fmt.Sprintf("  %v%v%v", keyv, f.Sep, f.Options[key]))
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

func (f *FOptionsMap) Validate(data []byte) error {
	ds := string(data)
	if f.CaseSensitive {
		for k := range f.Options {
			if k == ds {
				return nil
			}
		}
	} else {
		ds = strings.ToLower(ds)
		for k := range f.Options {
			if strings.ToLower(k) == ds {
				return nil
			}
		}
	}
	return errors.New("Invalid option.")
}

func (f *FOptionsMap) DefaultValue() string {
	return f.Default
}
