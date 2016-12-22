# Go semaphore

Implements a semaphore in Go, allowing to limit the concurrency of goroutines.

Note that this builds on top of Go channels, which is a recommended pattern used inside stdlib, but may be somewhat heavyweight.

Installation:

    go get github.com/andreyvit/sem

[Docs on godoc.org](https://godoc.org/github.com/andreyvit/sem)


## Example

```go
import (
    "github.com/andreyvit/sem"
)

func main() {
    s := sem.New(2)
    var wg sync.WaitGroup

    for i := 0; i < 30; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for k := 0; k < 10; k++ {
                s.Exec(func() {
                    fetchSomething()
                })
            }
        }()
    }

    wg.Wait()
    fmt.Printf("Finished %v requests.", atomic.LoadInt32(&requests))
}
```


## Usage

You can run a function without a return value:

```go
s.Exec(func() {
    fetchSomething()
})
```

or a function with an error return value:

```go
err := s.Exec(func() error {
    return fetchSomething()
})
```

or manually acquire and release:

```go
func fetchSomethingWithLimit(s sem.Sem) {
    s.Acquire()
    defer s.Release()
    // do something here
}
```
