/*
Package testingutil provides utilities not found in the core testing library.

Unfortunately, because these are functions and not methods of testing, the error
coordinates given by these utilities are less clean than they could have been.
*/
package testingutil

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

// ExpectEqual tests actual and expected for equality using reflect.DeepEqual and calls
// t.Errorf with a formatted error if the result is false.
//
// The optional desc parameters are formatted into the failure message if one is needed, and may
// either be of type string, ...interface{}, in which fmt.Sprintf is used to format them further,
// or anything else, in which case they are processed with fmt.Sprint.
func ExpectEqual(t *testing.T, actual, expected interface{}, desc ...interface{}) {
	if reflect.DeepEqual(actual, expected) {
		return
	}

	_, file, line, _ := runtime.Caller(1)
	desc1 := fmt.Sprintf("%s:%d%s", file, line, formatDesc(desc))
	t.Errorf("%s\nActual:   %#v\nExpected: %#v\n", desc1, actual, expected)
}

// ExpectDie runs f and tests whether it caused a panic. If not, it calls t.Errorf.
//
// The same formatting convention is used as in ExpectEqual.
func ExpectDie(t *testing.T, f func(), desc ...string) {
	_, file, line, _ := runtime.Caller(1)
	defer func() {
		if x := recover(); x == nil {
			t.Errorf("%s:%d%s\nExpected panic", file, line, formatDesc(desc))
		}
	}()
	f()
}

// formatDesc formats desc using fmt.Sprintf if the first element is a string,
// or with fmt.Sprint otherwise.
func formatDesc(desc ...interface{}) string {
	switch len(desc) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(" %v", desc[0])
	default:
		if f, ok := desc[0].(string); ok {
			return fmt.Sprintf(" "+f, desc[1:])
		} else {
			return " " + fmt.Sprint(desc...)
		}
	}
	panic("not reached")
}
