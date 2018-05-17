package cleanup

// Func is a cleanup function. It may not be nil.
//
// A Func should not panic during normal operation.
type Func = func() error

func assert(f Func) {
	if f == nil {
		panic("The cleanup func must not be nil")
	}
}

// Recover wraps f to recover panics.
//
// The Transform t is only applied to a recovered error.
//
// If the panic value does not satisfy the error interface, it is repaniced.
func Recover(t Transform, f Func) Func {
	assert(f)
	return func() (err error) {
		defer func() {
			if x := recover(); x != nil {
				if e, ok := x.(error); ok {
					err = apply(t, e)
				} else {
					panic(x)
				}
			}
		}()

		return f()
	}
}

// Compose returns a Func equivalent to t(f()).
func Compose(f Func, t Transform) Func {
	assert(f)
	if t == nil {
		return f
	}
	return func() error {
		return t(f())
	}
}
