package cleanup

// A Transform is function that may or may not alter an error.
// A nil Transform is valid, and equivalent to be the identity Transform.
// A nonnil Transform is always called, even if the err is nil.
type Transform = func(error) error

func apply(t Transform, err error) error {
	if t == nil {
		return err
	}
	return t(err)
}

// Transformers returns a Transform that runs all Transforms on err.
func Transformers(ts ...Transform) Transform {
	if len(ts) == 0 {
		return nil
	}
	return func(err error) error {
		for _, t := range ts {
			err = apply(t, err)
		}
		return err
	}
}

// Note returns an identity transform that calls f if err != nil.
func Note(f func(error)) Transform {
	if f == nil {
		return nil
	}
	return func(err error) error {
		if err != nil {
			f(err)
		}
		return err
	}
}

// IfErr returns a Transform that only calls t if err != nil.
func IfErr(t Transform) Transform {
	if t == nil {
		return nil
	}
	return func(err error) error {
		if err != nil {
			return t(err)
		}
		return nil
	}
}
