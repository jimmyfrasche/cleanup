package cleanup_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/jimmyfrasche/cleanup"
)

// A Transform may ignore an error.
func ExampleTransform_ignoreEOF() {
	var t cleanup.Transform = func(err error) error {
		if err == io.EOF {
			return nil
		}
		return err
	}

	fmt.Println(t(nil))
	fmt.Println(t(io.EOF))
	fmt.Println(t(errors.New("other")))

	// Output:
	// <nil>
	// <nil>
	// other
}

// A transform may create an error out of nothing.
func ExampleTransform_expectedError() {
	var t cleanup.Transform = func(err error) error {
		if err == nil {
			return errors.New("expected error but got nil")
		}
		return err
	}

	fmt.Println(t(errors.New("error")))
	fmt.Println(t(nil))

	// Output:
	// error
	// expected error but got nil
}

// A Transform may make note of an error without doing anything to it.
func ExampleNote() {
	t := cleanup.Note(func(err error) {
		fmt.Print(err)
	})

	t(nil)
	t(errors.New("note me"))

	// Output:
	// note me
}

// A Transform may wrap an error.
func ExampleIfErr() {
	type myErr struct {
		error
	}
	t := cleanup.IfErr(func(err error) error {
		return &myErr{err}
	})

	fmt.Printf("%T\n", t(nil))
	fmt.Printf("%T\n", t(errors.New("nonnil")))

	// Output:
	// <nil>
	// *cleanup_test.myErr
}
