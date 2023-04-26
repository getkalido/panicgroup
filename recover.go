package panicgroup

import (
	"runtime/debug"

	"github.com/pkg/errors"
)

// ErrRecover converts a recover() call into an error. Prefer using
// WrapEgGoWithRecover.
//
// The caller must call recover(), because the stack changes when ErrRecover
// is called.
//
// See https://play.golang.org/p/FYD2QlmHHO9
//
// Usage
// 	defer func() {
// 		err = helpers.ErrRecover(recover())
// 	}()
func ErrRecover(errI interface{}) error {
	if errI != nil {
		switch v := errI.(type) {
		case error:
			return errors.Wrap(v, string(debug.Stack()))
		default:
			return errors.Errorf("recovered %v stack: %s", v, debug.Stack())
		}
	}
	return nil
}

// WrapEgGoWithRecover takes an input function f and wraps it in panic
// recovery code. If a panic occurs, it is captured and instead
// converted to a regular error.
func WrapEgGoWithRecover(f func() error) func() error {
	return func() (err error) {
		defer func() {
			errR := ErrRecover(recover())
			// Only override the error returned if the error recovered is not nil
			if errR != nil {
				err = errR
			}
		}()

		return f()
	}
}
