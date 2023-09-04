package form

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestEncodePrimitiveValue(t *testing.T) {
	a := assert.New(t)

	testEncode(a, true, "=true")
	testEncode(a, false, "=false")
	testEncode(a, 0, "=0")
	testEncode(a, 1, "=1")
	testEncode(a, -1, "=-1")
	testEncode(a, uint(0), "=0")
	testEncode(a, 1.0, "=1")
	testEncode(a, -1.50, "=-1.5")
	testEncode(a, 1.0/3, "=0.3333333333333333")
	testEncode(a, 2.0/3, "=0.6666666666666666")
	testEncode(a, "", "=")
	testEncode(a, "test", "=test")
	testEncode(a, "Hello world", "=Hello+world")
}

func TestUnsupportedType(t *testing.T) {
	a := assert.New(t)

	_, err := Marshal(nil)
	a.NotNilNow(err)

	_, err = Marshal(complex(1, 1))
	a.NotNilNow(err)
}

func testEncode(a *assert.Assertion, v any, expect string) {
	data, err := Marshal(v)
	a.NilNow(err)
	a.EqualNow(string(data), expect)
}
