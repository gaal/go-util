package regexputil

import (
	"fmt"
	"github.com/gaal/go-util/testingutil"
	"testing"

	"regexp"
)

func TestExtractSubmatch(t *testing.T) {
	input := []byte("foo:bar:42")
	re := regexp.MustCompile(`(\w+):(\w+):(\d+)`)

	var s string
	var b []byte
	var i int
	testingutil.ExpectEqual(
		t,
		ExtractSubmatch(re, input, &s, &b, &i),
		nil,
		"ExtractSubmatch")
	testingutil.ExpectEqual(
		t, input, []byte("foo:bar:42"), "ExtractSubmatch does not modify input")
	testingutil.ExpectEqual(t, s, "foo", "string extraction")
	testingutil.ExpectEqual(t, b, []byte("bar"), "[]byte extraction")
	testingutil.ExpectEqual(t, i, 42, "int extraction")

	b[2]++
	testingutil.ExpectEqual(t, b, []byte("bas"), "[]byte is writable")
	testingutil.ExpectEqual(t, input, []byte("foo:bas:42"), "[]byte into input")
}

func ExampleExtractSubmatch() {
	input := []byte("foo:bar:42")
	re := regexp.MustCompile(`(\w+):(\w+):(\d+)`)

	var s string
	var b []byte
	var i int
	if err := ExtractSubmatch(re, input, &s, &b, &i); err != nil {
		panic(fmt.Errorf("expected match: %v", err))
	}
	fmt.Println(s, b, i)
	// Output:
	// foo [98 97 114] 42
}
