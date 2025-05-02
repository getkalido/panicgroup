package panicgroup

// Wrap0 wraps a function which takes no arguments in panic recovery code.
func Wrap0(f func() error) func() error {
	return func() (err error) {
		defer func() {
			errR := ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		return f()
	}
}

// Wrap1 wraps a function which takes one argument in panic recovery code.
func Wrap1[T any, F func(T) error](f F) F {
	return func(t T) (err error) {
		defer func() {
			errR := ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		return f(t)
	}
}

// Wrap2 wraps a function which takes two arguments in panic recovery code.
func Wrap2[T1 any, T2 any, F func(T1, T2) error](f F) F {
	return func(t1 T1, t2 T2) (err error) {
		defer func() {
			errR := ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		return f(t1, t2)
	}
}

// WrapReturn wraps a function which takes an argument and returns one value
// (in addition to an error) in panic recovery code.
func WrapReturn[T any, R any, F func(T) (R, error)](f F) F {
	return func(t T) (result R, err error) {
		defer func() {
			errR := ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		return f(t)
	}
}

// WrapReturn wraps a function which takes two arguments and returns one value
// (in addition to an error) in panic recovery code.
func WrapReturn2[T1 any, T2 any, R any, F func(T1, T2) (R, error)](f F) F {
	return func(t1 T1, t2 T2) (result R, err error) {
		defer func() {
			errR := ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		return f(t1, t2)
	}
}
