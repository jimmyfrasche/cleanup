package cleanup_test

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jimmyfrasche/cleanup"
)

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("wrapped: %v", err)
}

func f() (err error) {
	var h cleanup.Handler

	// defer h.Cleanup(nil) is also valid.
	// This construction lets us handle returned errors
	// without doing anything special at the call site.
	defer func() {
		// If you only care about the first error, call First instead of Err.
		err = h.Cleanup(err).Err()
	}()

	// Wrap all non-nil errors during Cleanup.
	h.Transform(Wrap)

	f, err := ioutil.TempFile("", "")
	if err != nil {
		// If this err is nonnil, it will still pass through Wrap.
		return err
	}

	h.Defer(func() error {
		return os.Remove(f.Name())
	})

	// Close will be called before os.Remove above.
	h.Defer(f.Close)

	h.Defer(func() error {
		return errors.New("simulated teardown failure 1")
	})

	w := bufio.NewWriter(f)

	h.Defer(func() error {
		// Don't bother flushing if there are errors.
		// I'm not sure why you would do that,
		// but coming up with a good example is hard.
		if len(h.Errors) > 0 {
			return errors.New("skipping Flush")
		}
		return w.Flush()
	})

	h.Defer(func() error {
		return errors.New("simulated teardown failure 2")
	})

	_, err = w.WriteString("do some work")

	// This is redundant given the defer above, but perfectly safe.
	// Something like this might be preferred to keep the defer simpler
	// or to have separate logic for errors during setup vs. teardown.
	return h.Cleanup(err).Err()
}

func ExampleHandler() {
	fmt.Println(f())

	// Output:
	// wrapped: simulated teardown failure 2
	// wrapped: skipping Flush
	// wrapped: simulated teardown failure 1
}
