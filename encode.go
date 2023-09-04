// package go-form implements encoding of URL encoded form data.
package form

import (
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Marshal returns the URL encoded form data of v.
func Marshal(v any) ([]byte, error) {
	e := &encoderState{
		buf: map[string][]string{},
	}

	if mv, ok := v.(map[string][]string); ok {
		e.buf = mv
		return e.encode(), nil
	}

	if err := valueEncoder(e, reflect.ValueOf(v), ""); err != nil {
		return nil, err
	}

	return e.encode(), nil
}

type encoderState struct {
	buf map[string][]string
}

func (e *encoderState) appendToKey(k, v string) {
	_, ok := e.buf[k]
	if ok {
		e.buf[k] = append(e.buf[k], v)
	} else {
		e.buf[k] = []string{v}
	}
}

func (e *encoderState) encode() []byte {
	keys := make([]string, 0, len(e.buf))
	for k := range e.buf {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := strings.Builder{}

	for i, k := range keys {
		if i > 0 {
			buf.WriteRune('&')
		}

		v := e.buf[k]
		if len(v) == 1 {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteRune('=')
			buf.WriteString(url.QueryEscape(v[0]))
		} else {
			for vi, vv := range v {
				if i > 0 {
					buf.WriteRune('&')
				}

				buf.WriteString(url.QueryEscape(k))
				buf.WriteRune('[')
				buf.WriteString(strconv.Itoa(vi))
				buf.WriteRune(']')

				buf.WriteRune('=')

				buf.WriteString(url.QueryEscape(vv))
			}
		}
	}

	return []byte(buf.String())
}

func valueEncoder(e *encoderState, v reflect.Value, prefix string) error {
	switch v.Kind() {
	case reflect.Bool:
		boolEncoder(e, v, prefix)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intEncoder(e, v, prefix)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		uintEncoder(e, v, prefix)
	case reflect.Float32, reflect.Float64:
		floatEncoder(e, v, prefix)
	case reflect.Array, reflect.Slice:
		// TODO
	case reflect.Map:
		// TODO
	case reflect.String:
		stringEncoder(e, v, prefix)
	case reflect.Struct:
		// TODO
	case reflect.Pointer:
		valueEncoder(e, v.Elem(), prefix)
	default:
		return ErrUnsupportedType
	}

	return nil
}

func boolEncoder(e *encoderState, v reflect.Value, prefix string) {
	if v.Bool() {
		e.appendToKey(prefix, "true")
	} else {
		e.appendToKey(prefix, "false")
	}
}

func intEncoder(e *encoderState, v reflect.Value, prefix string) {
	e.appendToKey(prefix, strconv.FormatInt(v.Int(), 10))
}

func uintEncoder(e *encoderState, v reflect.Value, prefix string) {
	e.appendToKey(prefix, strconv.FormatUint(v.Uint(), 10))
}

func floatEncoder(e *encoderState, v reflect.Value, prefix string) {
	e.appendToKey(prefix, strconv.FormatFloat(v.Float(), 'f', -1, 64))
}

func stringEncoder(e *encoderState, v reflect.Value, prefix string) {
	e.appendToKey(prefix, v.String())
}
