package parmap

import (
	"context"
	"runtime"
	"testing"
)

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

const N, M = 30, 31

func BenchmarkNoParmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(fib(N))
		runtime.KeepAlive(fib(M))
	}
}

func BenchmarkParmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		type result struct {
			a, b int
		}
		r, err := New[result](context.Background()).
			F(func(r *result) {
				r.a = fib(N)
			}).
			FErr(func(r *result) error {
				r.b = fib(M)
				return nil
			}).
			Done()
		runtime.KeepAlive(r)
		runtime.KeepAlive(err)
	}
}
