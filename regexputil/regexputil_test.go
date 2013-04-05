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

func TestReplaceFirst(t *testing.T) {
	testFirstReplacements := []struct {
		search         string
		replace        string
		testInput      string
		expectedOutput string
	}{
		// search, replace, input, output
		{"a", "b", "a", "b"},
		{"a", "b", "c", "c"},
		{"(a)", "b", "a", "b"},
		{"b(a)", "$1", "ba", "a"},
		{"(a)", "$1$1$1", "a", "aaa"},
		{"(a)", "$1$1$1", "bab", "baaab"},
		{"(.)(.)", "$2$1", "ab", "ba"},
		{`(\d)(\d)`, "$2$1", "a12b", "a21b"},
		{`(\d)(\d)`, "$2$1", "a12b12c", "a21b12c"},
	}

	for i, rep := range testFirstReplacements {
		testingutil.ExpectEqual(
			t,
			ReplaceFirst(regexp.MustCompile(rep.search), []byte(rep.testInput), []byte(rep.replace)),
			[]byte(rep.expectedOutput),
			"[%d] ReplaceFirst(%q, %q, %q)", i, rep.search, rep.testInput, rep.replace)

		testingutil.ExpectEqual(
			t,
			ReplaceFirstString(regexp.MustCompile(rep.search), rep.testInput, rep.replace),
			rep.expectedOutput,
			"[%d] ReplaceFirstString(%q, %q, %q)", i, rep.search, rep.testInput, rep.replace)

	}
}

func ExampleReplaceFirstString() {
	fmt.Println(ReplaceFirstString(regexp.MustCompile(`(\d)(\d)`), "a12b12c", "$2$1"))
	// Output:
	// a21b12c
}
