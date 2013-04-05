package testingutil

import (
	"testing"
)

func TestExpectEqual(t *testing.T) {
	p1 := true
	p2 := true
	deep := []struct {
		aa int
		bb string
		pp *bool
	}{
		{42, "hello", &p1},
		{42, "hello", &p2},
	}
	ExpectEqual(t, deep[0], deep[1], "deep")

	want := "hello"
	ExpectEqual(t, "hello", want, "string")

	bad := []interface{}{
		[]byte("hello"),
		"helo",
		42.2,
		nil,
	}
	for i, v := range bad {
		t1 := new(testing.T)
		ExpectEqual(t1, v, want, "expected mismatch")
		ExpectEqual(t, t1.Failed(), true, "[%d]: ExpectEqual on expected mismatch", i)
	}
}

func TestExpectDie(t *testing.T) {
	ExpectDie(t, func() { panic("aaaaahh") }, "simple panic dies")

	t1 := new(testing.T)
	ExpectDie(t1, func() {}, "doesn't die")
	ExpectEqual(t, t1.Failed(), true, "ExpectDie on something that doesn't die fails")
}

func ExampleExpectDie() {
	// In your test, this will be passed in as a paramter.
	t := new(testing.T)

	ExpectDie(t, func() { /* I will not die! */
	}, "stubborn code does not die")
	// This will fail with a message like:
	// testingutil.go:39:      /Users/me/src/mycode_test.go:48 [stubborn code does not die]
	//       Expected panic
	//       FAIL
	// FAIL    github.com/gaal/go-util/testingutil     0.010s
}

func ExampleExpectEqual() {
	// In your test, this will be passed in as a paramter.
	t := new(testing.T)

	want := "hello"
	ExpectEqual(t, "hello", want, "expected match")
	ExpectEqual(t, "goodbye", want, "goodbye == hello?")
	// This will fail with a message like:
	// testingutil.go:29:      /Users/me/src/mycode_test.go:51: [goodbye == hello?]
	//        Actual: "goodbye"
	//        Expected:   "hello"
	//        FAIL
	// FAIL    mycode     0.010s

	// structs can also be passed for comparison.
}
