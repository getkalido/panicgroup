package panicgroup_test

import (
	"errors"
	"fmt"

	"github.com/getkalido/panicgroup"
)

var ErrSomethingTerrible error = errors.New("something terrible has happened!")

// ErrorWrapping illustrates how ErrRecover can be used to recover from
// a panic and convert it into a useful error instead. This allows
// the caller of the panicking function to regain control and handle
// the panic as desired.
//
// WrapEgGoWithRecover is simply a wrapper to avoid needing to manually
// do this wrapping.
func ExampleErrRecover_errorWrapping() {
	// An example function, which has a chance to panic.
	panicFunc := func() (err error) {
		defer func() {
			// ErrRecover should recover from the panic and convert it to
			// an error.
			errR := panicgroup.ErrRecover(recover())
			if errR != nil {
				err = errR
			}
		}()
		panic(ErrSomethingTerrible)
	}

	err := panicFunc()
	if err == nil {
		fmt.Println("No error was captured.")
	} else {
		fmt.Println("The panic was successfully converted into an error.")
	}
	// Output:
	// The panic was successfully converted into an error.
}

// SimplifiedWrapping illustrates how WrapEgGoWithRecover can be used
// to recover from a panic and convert it into a useful error instead.
// This allows the caller of the panicking function to regain control
// and handle the panic as desired.
//
// This simplifies the manual process of adding panic recovery to a
// block of code.
func ExampleWrapEgGoWithRecover_simplifiedWrapping() {
	// An example function, which has a chance to panic.
	panicFunc := panicgroup.WrapEgGoWithRecover(func() (err error) {
		panic(ErrSomethingTerrible)
	})

	err := panicFunc()
	if err == nil {
		fmt.Println("No error was captured.")
	} else {
		fmt.Println("The panic was successfully converted into an error.")
	}
	// Output:
	// The panic was successfully converted into an error.
}
