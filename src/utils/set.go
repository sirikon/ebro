package utils

import "slices"

type Set[T comparable] struct {
	s []T
	m map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{s: []T{}, m: map[T]bool{}}
}

func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		if _, ok := s.m[item]; !ok {
			s.s = append(s.s, item)
			s.m[item] = true
		}
	}
}

func (s *Set[T]) Delete(items ...T) {
	for _, item := range items {
		if _, ok := s.m[item]; ok {
			pos := slices.Index(s.s, item)
			s.s = slices.Delete(s.s, pos, pos+1)
			delete(s.m, item)
		}
	}
}

func (s *Set[T]) Get(pos int) T {
	return s.s[pos]
}

func (s *Set[T]) Length() int {
	return len(s.s)
}

func (s *Set[T]) List() []T {
	return slices.Clone(s.s)
}
