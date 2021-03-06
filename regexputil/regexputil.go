/*
Package regexputil provides utilities not found in the core regexp library.
*/
package regexputil

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	ErrCount = errors.New("wrong value count")
	ErrMatch = errors.New("no match")
	ErrType  = errors.New("wrong type")
)

// ExtractSubmatch performs a regexp match on src, and attempts to extract all 
// submatches into values, which must be of pointer type.
//
// If the match fails, or if len(values) != re.NumSubexp(), an error is returned
// and values are not modified. If the match succeeds, but any of the values cannot
// be written to, an error is returned and values are not defined.
//
// Supported types for values:
//   *[]byte - The target byte slice as returned by regexp. May be writable.
//   *string - Does a string conversion, and therefore a copy.
//   *int    - Converted with strconv.Atoi. For other numeric conversions,
//             extract to a string and perform one yourself.
//
// In the future, addiional numeric types and some scanning interface may be contemplated.
func ExtractSubmatch(re *regexp.Regexp, src []byte, values ...interface{}) error {
	if re.NumSubexp() != len(values) {
		return ErrCount
	}
	sm := re.FindSubmatch(src)
	if sm == nil {
		return ErrMatch
	}
	for i, val := range values {
		if err := extractTo(val, sm[i+1]); err != nil {
			return err
		}
	}
	return nil
}

func extractTo(val interface{}, src []byte) error {
	switch v := val.(type) {
	case *[]byte:
		*v = src
	case *string:
		*v = string(src)
	case *int:
		i, err := strconv.Atoi(string(src))
		if err != nil {
			return err
		}
		*v = i
	default:
		return ErrType
	}
	return nil
}

// ReplaceFirst returns a copy of src, replacing the first match of the Regexp with the replacement
// text repl. Inside repl, $ signs are interpreted as in regexp.Expand, so for instance $1
// represents the text of the first submatch.
//
// This function complements regexp.ReplaceAll.
func ReplaceFirst(re *regexp.Regexp, src, repl []byte) []byte {
	if m := re.FindSubmatchIndex(src); m != nil {
		out := make([]byte, m[0])
		copy(out, src[0:m[0]])
		out = re.Expand(out, repl, src, m)
		if m[1] < len(src) {
			out = append(out, src[m[1]:]...)
		}
		return out
	}
	out := make([]byte, len(src))
	copy(out, src)
	return out
}

// ReplaceFirstString is like ReplaceFirst, except that src and repl are strings, as is the
// return value.
//
// This function complements regexp.ReplaceAllString.
func ReplaceFirstString(re *regexp.Regexp, src, repl string) string {
	return string(ReplaceFirst(re, []byte(src), []byte(repl)))
}
