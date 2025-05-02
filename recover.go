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
// Usage:
//
//	defer func() {
//		errR := helpers.ErrRecover(recover())
//		if errR != nil {
//			err = errR
//		}
//	}()
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
	return WrapEgGoWithCustomRecover(f, func(err error) error { return err })
}

// ErrorHandler represents a function which performs some action on an error
// which was converted from a panic.
//
// This function should return as soon as possible if an error group was
// canceled.
type ErrorHandler func(error) error

// WrapEgGoWithCustomRecover takes an input function f and wraps it in panic
// recovery code. If a panic occurs, it is captured and instead
// converted to a regular error. This error is then passed to recovery, so that
// recovery can decide what happens to the error.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	errChan := make(chan error, 1)
//	wrappedFunc := panicgroup.WrapEgGoWithCustomRecover(
//		func() {
//			panic("AAAAAAAH!")
//		},
//		func(err error) error {
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			case errChan<-err:
//				return err
//			}
//		}
//	)
func WrapEgGoWithCustomRecover(f func() error, recovery ErrorHandler) func() error {
	return func() (err error) {
		defer func() {
			errR := ErrRecover(recover())
			// Only override the error returned if the error recovered is not nil
			if errR != nil {
				err = recovery(errR)
			}
		}()

		return f()
	}
}
