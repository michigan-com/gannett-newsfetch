// Package sem implements a semaphore, allowing to limit the concurrency of goroutines.
package sem

type Sem chan struct{}

// New makes a new semaphore with a given limit on concurrency.
func New(n int) Sem {
	return make(chan struct{}, n)
}

func (s Sem) Cap() int {
	return cap(s)
}

// Acquire marks the start of a limited concurrency section.
func (s Sem) Acquire() {
	s <- struct{}{}
}

// Release marks the end of a limited concurrency section.
func (s Sem) Release() {
	<-s
}

// Executes the given function with a limit on concurrency.
func (s Sem) Exec(f func()) {
	s.Acquire()
	defer s.Release()
	f()
}

// Executes the given function with a limit on concurrency, returning an error.
func (s Sem) Exece(f func() error) error {
	s.Acquire()
	defer s.Release()
	return f()
}

// Executes the given function with a limit on concurrency, returning a function.
func (s Sem) Execfun(f func() func()) func() {
	s.Acquire()
	defer s.Release()
	return f()
}

// Acquires all the spots on the semaphore, effectively making sure that no other functions can execute.
func (s Sem) Drain() {
	n := cap(s)
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
}
