package regexputil

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractSubmatch(t *testing.T) {
	input := []byte("foo:bar:42")
	re := regexp.MustCompile(`(\w+):(\w+):(\d+)`)

	type res struct {
		S string
		B []byte
		I int
	}
	var r res
	if got := ExtractSubmatch(re, input, &r.S, &r.B, &r.I); got != nil {
		t.Errorf("ExtractSubmatch(%q, %q, *string, *[]byte, *int) = %q, want nil", re, input, got)
	}
	if !bytes.Equal(input, []byte("foo:bar:42")) {
		t.Errorf("ExtractSubmatch should not modify input, got: %s", input)
	}
	want := res{"foo", []byte("bar"), 42}
	if diff := cmp.Diff(want, r); diff != "" {
		t.Errorf("ExtractSubmatch(%q, %q, *string, *[]byte, *int): diff(-want,+got)\n%s", re, input, diff)
	}
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
	data := []struct {
		search    string
		replace   string
		testInput string
		want      string
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

	for _, d := range data {
		re := regexp.MustCompile(d.search)
		if got := ReplaceFirst(re, []byte(d.testInput), []byte(d.replace)); string(got) != d.want {
			t.Errorf("ReplaceFirst(%q, %q, %q) = %q, want = %q", d.search, d.testInput, d.replace, string(got), d.want)
		}
		if got := ReplaceFirstString(re, d.testInput, d.replace); got != d.want {
			t.Errorf("ReplaceFirstString(%q, %q, %q) = %q, want = %q", d.search, d.testInput, d.replace, got, d.want)
		}
	}
}

func ExampleReplaceFirstString() {
	fmt.Println(ReplaceFirstString(regexp.MustCompile(`(\d)(\d)`), "a12b12c", "$2$1"))
	// Output:
	// a21b12c
}
