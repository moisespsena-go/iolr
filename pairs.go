package iou

import (
	"fmt"

	"go4.org/sort"
)

type StringPair struct {
	K, V string
	Ki   interface{}
}

func StringPairSortByValue(a, b *StringPair) bool {
	return a.V < b.V
}

func StringPairSortByKey(a, b *StringPair) bool {
	return a.K < b.K
}

type StringPairs []*StringPair

func (s StringPairs) Sort(less func(i, j *StringPair) bool) StringPairs {
	sort.Slice(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
	return s
}

func (s StringPairs) Append(pairs ...*StringPair) StringPairs {
	s = append(s, pairs...)
	return s
}

func (s StringPairs) AddPairs(kv ...interface{}) StringPairs {
	if len(kv)%2 == 1 {
		panic("invalid pairs count")
	}
	for i, l := 0, len(kv); i < l; i += 2 {
		s = s.Add(kv[i], kv[i+1].(string))
	}
	return s
}

func (s StringPairs) Add(key interface{}, value string) StringPairs {
	var ks string
	if value == "" {
		ks = ""
		key = ""
	} else if k, ok := key.(string); ok {
		ks = k
	} else {
		ks = fmt.Sprint(key)
	}
	return s.Append(&StringPair{ks, value, key})
}

func (s StringPairs) AddTitle(title string) StringPairs {
	return s.Add("", title)
}

func (s StringPairs) AddBlank() StringPairs {
	return s.Add("", "")
}

func StringsToPairs(v ...string) (pairs StringPairs) {
	for i, v := range v {
		pairs = pairs.Add(i+1, v)
	}
	return
}

func NewPairs(kv ...interface{}) (pairs StringPairs) {
	return pairs.AddPairs(kv...)
}

func MapToPairs(m map[interface{}]string) (pairs StringPairs) {
	for k, v := range m {
		pairs = pairs.Add(k, v)
	}
	return pairs
}
