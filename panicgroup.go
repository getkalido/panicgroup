// Package errgroup provides utility functions and types for easily
// recovering from panics.
package panicgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// WgErrGroup re-implements the core functionality of [errgroup.Group], but
// wrapped in a panic-handling function.
type WgErrGroup interface {
	Go(func() error)
	Wait() error
}

// WgErrGroupWithCustom re-implements the core functionality of
// [errgroup.Group], but wrapped in a panic-handling function. It also adds
// a way of specifying a custom recovery function, which will execute if
// a panic is recovered.
type WgErrGroupWithCustom interface {
	WgErrGroup
	GoCustomRecover(f func() error, errorHandler func(error) error)
}

type wgErrGroup = Group

type Group struct {
	wg *errgroup.Group
}

var _ WgErrGroup = (*Group)(nil)
var _ WgErrGroupWithCustom = (*Group)(nil)

// WaitGroup returns a new instance of WgErrGroup.
//
// WaitGroupWithContext should be preferred, since it is able to cancel
// execution if one of the goroutines started with Go encounters an error.
func WaitGroup() *Group {
	return &Group{wg: &errgroup.Group{}}
}

// WithContext returns a new WgErrGroup and an associated Context
// derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	errGroup, errContext := errgroup.WithContext(ctx)
	return &Group{wg: errGroup}, errContext
}

// WaitGroupWithContext returns a new WgErrGroup and an associated Context
// derived from ctx.
//
// This function is provided purely for backwards compatibility.
//
// See [WithContext] for the real implementation.
func WaitGroupWithContext(ctx context.Context) (*Group, context.Context) {
	return WithContext(ctx)
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
//
// If the captured panic needs custom handling, use [Group.GoCustomRecover]
// instead.
func (g *Group) Go(f func() error) {
	g.wg.Go(WrapEgGoWithRecover(f))
}

// GoCustomRecover calls the given function in a new goroutine. If a panic
// occurs, the value will be recovered and wrapped with [ErrRecover]. The
// resulting error is then passed to recovery and the error returned from
// recovery will be returned to the [Group].
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
//
// If recovery contains blocking operations like sending on a channel,
// recovery must return when the [Group]'s context is canceled.
//
// If the captured panic just needs to be returned, use [Group.Go] instead.
func (g *Group) GoCustomRecover(f func() error, recovery func(error) error) {
	g.wg.Go(WrapEgGoWithCustomRecover(f, recovery))
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	return g.wg.Wait()
}
