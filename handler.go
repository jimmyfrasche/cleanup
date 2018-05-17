package cleanup

import "github.com/jimmyfrasche/multiple"

// Handler manages sequences of cleanup functions.
//
// A Handler may not be reused and is not safe for concurrent access.
type Handler struct {
	ran bool
	fs  []Func
	t   Transform

	Errors multiple.Errors
}

// Defer a Func to run.
//
// It is an error to add a func after Cleanup has been called.
func (h *Handler) Defer(f Func) {
	if h.ran {
		panic("attempt to Defer after Cleanup")
	}
	assert(f)
	h.fs = append(h.fs, f)
}

// Transform sets a global Transform applied to every Func during Cleanup.
func (h *Handler) Transform(t Transform) {
	if h.ran {
		panic("attempt to set Transform after Cleanup")
	}
	h.t = t
}

// Cleanup runs all registered handlers in reverse order
// and collects any errors in Errors.
//
// If err is not nil, it will be transformed and appended to the error stack.
//
// After the first call, subsequent calls will be ignored.
func (h *Handler) Cleanup(err error) multiple.Errors {
	if h.ran {
		return h.Errors
	}
	h.ran = true

	_ = h.Errors.Append(apply(h.t, err))

	for i := len(h.fs) - 1; i >= 0; i-- {
		_ = h.Errors.Append(apply(h.t, h.fs[i]()))
	}

	return h.Errors
}
