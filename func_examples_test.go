package cleanup_test

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/jimmyfrasche/cleanup"
)

func ExampleRecover() {
	f := cleanup.Recover(nil, func() error {
		panic(errors.New("panic"))
		return errors.New("return")
	})

	fmt.Println(f())

	// Output:
	// panic
}

func ExampleRecover_transform() {
	// This transform re-panics runtime errors but lets all others through as-is.
	Transform := func(err error) error {
		if rerr, ok := err.(runtime.Error); ok {
			panic(rerr)
		}
		return err
	}

	// This does not cause a runtime error so it is handled normally.
	g := cleanup.Recover(Transform, func() error {
		panic(errors.New("caught panic"))
		return nil
	})
	fmt.Println(g())

	// This causes a runtime error which our Transform specifically re-panics.
	f := cleanup.Recover(Transform, func() error {
		a, b := 1, 0
		_ = a / b
		return nil
	})

	defer func() {
		fmt.Println(recover())
	}()

	f()

	// Output:
	// caught panic
	// runtime error: integer divide by zero
}

func ExampleCompose() {
	log := cleanup.Note(func(err error) {
		fmt.Println("log:", err)
	})

	f := func() error {
		return nil
	}

	g := func() error {
		return errors.New("an error")
	}

	// This won't cause any log output.
	cleanup.Compose(f, log)()

	// This will.
	cleanup.Compose(g, log)()

	// Output:
	// log: an error
}
