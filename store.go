package stores

import (
	"errors"
	"sync"
)

type Store[T any] struct {
	mutex    sync.RWMutex
	contents storeEntries[T]
}

func (s *Store[T]) ensureEntries() {
	if s.contents == nil {
		s.contents = make(storeEntries[T])
	}
}

func (s *Store[T]) Set(key string, content *T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.ensureEntries()

	s.contents.set(key, content)
}

func (s *Store[T]) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.ensureEntries()

	s.contents.delete(key)
}

func (s *Store[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.ensureEntries()

	s.contents.clear()
}

func (s *Store[T]) Get(key string) (*T, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.ensureEntries()

	entry := s.contents.find(key)

	if entry == nil {
		return nil, errors.New("Element with the provided key does not exist.")
	}

	return entry, nil
}

func (s *Store[T]) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.ensureEntries()

	var keys []string

	for key := range s.contents {
		keys = append(keys, key)
	}

	return keys
}

func (s *Store[T]) Values() []*T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.ensureEntries()

	var values []*T

	for _, content := range s.contents {
		values = append(values, content)
	}

	return values
}

func (s *Store[T]) Entries() []Tuple[string, *T] {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.ensureEntries()

	var entries []Tuple[string, *T]

	for key, content := range s.contents {
		entries = append(entries, NewTuple(key, content))
	}

	return entries
}

func (s *Store[T]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.ensureEntries()

	return len(s.contents)
}
