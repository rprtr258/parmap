# Usage

```go
// define some type to store all results
type result struct {
    a int
    b int
    c int
}
r, err := parmap.New[result](ctx).
    F(func(r *result) {
        // use some pure functions to evaluate some fields
        r.a = fib(100)
    }).
    FErr(func(r *result) error {
        // or parse something
        b, err = parse()
        r.b = b
        return err
    }).
    FuncErr(func(ctx context.Context, r *result) error {
        // or extract data from external resource
        c, err := fetch(ctx, r.a)
        r.c = c
        return err
    })
    Done()
// all errors and panics are handled
if err != nil { ...  }
// r now contains all results which were obtained asynchronously
println(r)
```
