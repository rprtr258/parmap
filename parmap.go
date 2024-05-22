package parmap

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Parmap[T any] struct {
	grp    *errgroup.Group
	ctx    context.Context
	result T
}

func New[T any](ctx context.Context) *Parmap[T] {
	g, gctx := errgroup.WithContext(ctx)
	return &Parmap[T]{
		grp: g,
		ctx: gctx,
	}
}

func (p *Parmap[T]) Limit(n int) *Parmap[T] {
	p.grp.SetLimit(n)
	return p
}

func (p *Parmap[T]) FuncErr(f func(context.Context, *T) error) *Parmap[T] {
	p.grp.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				if err != nil {
					err = fmt.Errorf("err: %w, panic: %v", err, r)
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()
		return f(p.ctx, &p.result)
	})
	return p
}

func (p *Parmap[T]) Func(f func(context.Context, *T)) *Parmap[T] {
	return p.FuncErr(func(ctx context.Context, t *T) error {
		f(ctx, t)
		return nil
	})
}

func (p *Parmap[T]) FErr(f func(*T) error) *Parmap[T] {
	return p.FuncErr(func(ctx context.Context, t *T) error {
		return f(t)
	})
}

func (p *Parmap[T]) F(f func(*T)) *Parmap[T] {
	return p.FuncErr(func(ctx context.Context, t *T) error {
		f(t)
		return nil
	})
}
func (p *Parmap[T]) Done() (T, error) {
	p.grp.Wait()
	return p.result, nil
}
