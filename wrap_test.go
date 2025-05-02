package panicgroup_test

import (
	"fmt"

	"github.com/getkalido/panicgroup"
)

// BareWrapping illustrates how WrapReturn can be used to wrap functions
// in error recovery code. Note how the safe function acts normally and
// the panicking function returns an error.
//
// The same principles apply to the other wrapping functions.
func ExampleWrapReturn_bareWrapping() {
	// An example function, which has a chance to panic.
	panicFunc := panicgroup.WrapReturn(func(input string) (string, error) {
		panic(ErrSomethingTerrible)
	})

	result, err := panicFunc("this")
	if result != "" {
		fmt.Println("Panicked function returned a result.")
	}
	if err == nil {
		fmt.Println("No error was captured.")
	} else {
		fmt.Println("The panic was successfully converted into an error.")
	}

	safeFunc := panicgroup.WrapReturn(func(input string) (string, error) {
		return input + "!", nil
	})
	result, err = safeFunc("this")
	fmt.Printf("Safe function returned: %q.\n", result)
	if err != nil {
		fmt.Printf("The safe function returned an error: %q.\n", err.Error())
	}

	// Output:
	// The panic was successfully converted into an error.
	// Safe function returned: "this!".
}
