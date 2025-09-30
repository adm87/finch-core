package enum

import (
	"encoding/json"
	"fmt"

	"github.com/adm87/finch-core/linq"
)

type Enum[T ~int] interface {
	~int
	String() string
	IsValid() bool
}

func MarshalEnum[T Enum[T]](e T) ([]byte, error) {
	return json.Marshal(e.String())
}

func UnmarshalEnum[T Enum[T]](data []byte) (T, error) {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		var zero T
		return zero, err
	}
	return Value[T](str)
}

func Mapping[T Enum[T]]() map[string]T {
	m := make(map[string]T)
	for _, v := range Values[T]() {
		m[v.String()] = v
	}
	return m
}

func Names[T Enum[T]]() []string {
	return linq.Keys(Mapping[T]())
}

func Values[T Enum[T]]() []T {
	var vals []T
	for i := 0; ; i++ {
		e := T(i)
		if !e.IsValid() {
			break
		}
		vals = append(vals, e)
	}
	return vals
}

func Value[T Enum[T]](name string) (T, error) {
	if v, ok := Mapping[T]()[name]; ok {
		return v, nil
	}
	var zero T
	return zero, fmt.Errorf("invalid enum name: %s", name)
}
