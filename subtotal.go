package main

import (
	"fmt"
)

type Subtotal interface {
	Name() string
	HasParent() bool
}

type subtotal struct {
	name     string
	parentID int
}

func NewSubtotalWithNoParent(name string) (s *subtotal) {
	s = &subtotal{
		name: name,
	}

	return
}

func (s *subtotal) Name() string {
	return s.name
}

func (s *subtotal) HasParent() bool {
	return s.parentID != 0
}

type SubtotalAPI interface {
	Post(Subtotal) error
	GetByName(string) (Subtotal, error)
}

type inMemorySubtotalAPI struct {
	store map[string]Subtotal
}

func NewInMemorySubtotalAPI() (a *inMemorySubtotalAPI) {
	a = &inMemorySubtotalAPI{
		store: make(map[string]Subtotal),
	}

	return
}

func (a *inMemorySubtotalAPI) Post(s Subtotal) (e error) {
	var (
		exists bool
	)

	_, exists = a.store[s.Name()]
	if exists {
		e = fmt.Errorf(`Could not POST: Subtotal with name "%s" already exists`,
			s.Name(),
		)

		return
	}

	a.store[s.Name()] = s

	return
}

func (a *inMemorySubtotalAPI) GetByName(name string) (s Subtotal, e error) {
	var (
		exists bool
	)

	_, exists = a.store[name]
	if !exists {
		e = fmt.Errorf(`Could not GET by name: Subtotal "%s" does not exist`,
			name,
		)

		return
	}

	s = a.store[name]

	return
}
