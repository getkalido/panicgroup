// Package errgroup provides utility functions and types for easily
// recovering from panics.
package panicgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// WgErrGroup re-implements the core functionality of errgroup.Group, but
// wrapped in a panic-handling function.
type WgErrGroup interface {
	Go(func() error)
	Wait() error
}

type wgErrGroup struct {
	wg *errgroup.Group
}

// WaitGroup returns a new instance of WgErrGroup.
//
// WaitGroupWithContext should be preferred, since it is able to cancel
// execution if one of the goroutines started with Go encounters an error.
func WaitGroup() WgErrGroup {
	return &wgErrGroup{wg: &errgroup.Group{}}
}

// WaitGroupWithContext returns a new WgErrGroup and an associated Context
// derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WaitGroupWithContext(ctx context.Context) (WgErrGroup, context.Context) {
	errGroup, errContext := errgroup.WithContext(ctx)
	return &wgErrGroup{wg: errGroup}, errContext
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *wgErrGroup) Go(f func() error) {
	g.wg.Go(WrapEgGoWithRecover(f))
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *wgErrGroup) Wait() error {
	return g.wg.Wait()
}
