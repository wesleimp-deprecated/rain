package semerrgroup

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// Group is the Semphore ErrorGroup itself.
type Group interface {
	Go(func() error)
	Wait() error
}

// New returns a new Group of a given size.
func New(s int) Group {
	if s == 1 {
		return &serialGroup{}
	}

	return &parallelGroup{
		ch: make(chan bool, s),
		g:  errgroup.Group{},
	}
}

var _ Group = &parallelGroup{}

type parallelGroup struct {
	ch chan bool
	g  errgroup.Group
}

// Go execs one function respecting the group and semaphore.
func (s *parallelGroup) Go(fn func() error) {
	s.g.Go(func() error {
		s.ch <- true
		defer func() {
			<-s.ch
		}()
		return fn()
	})
}

// Wait waits for the group to complete and return an error if any.
func (s *parallelGroup) Wait() error {
	return s.g.Wait()
}

var _ Group = &serialGroup{}

type serialGroup struct {
	err     error
	errOnce sync.Once
}

// Go execs runs `fn` and saves the result if no error has been encountered.
func (s *serialGroup) Go(fn func() error) {
	if s.err != nil {
		return
	}
	if err := fn(); err != nil {
		s.errOnce.Do(func() {
			s.err = err
		})
	}
}

// Wait waits for Go to complete and returns the first error encountered.
func (s *serialGroup) Wait() error {
	return s.err
}
