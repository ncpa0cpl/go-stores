package stores

import (
	"sort"
	"time"
)

type CacheEntry[T any] struct {
	value       T
	accessCount int
	timestamp   int64
}

func (ce *CacheEntry[T]) evaluate() float64 {
	age := float64((time.Now().UnixMilli() - ce.timestamp) / 1000)

	if age <= 2 {
		return 100 + (float64(ce.accessCount+1) / age)
	}

	return float64(ce.accessCount+1) / age
}

type CacheStore[T any] struct {
	store          Store[CacheEntry[T]]
	entryLifetime  int
	maxSize        int
	cleanupCounter int
}

func (cs *CacheStore[T]) SetMaxStoreSize(size int) {
	cs.maxSize = size
	go cs.removeEntriesIfMaxSizeExceeded()
}

func (cs *CacheStore[T]) SetMaxEntryLifetime(ms int) {
	cs.entryLifetime = ms
}

type EntryTuples []Tuple[string, float64]

func (e EntryTuples) Len() int           { return len(e) }
func (e EntryTuples) Less(i, j int) bool { return e[i].b < e[j].b }
func (e EntryTuples) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (cs *CacheStore[T]) removeEntriesIfMaxSizeExceeded() {
	cs.cleanupCounter++
	currentCleanupCounter := cs.cleanupCounter

	time.Sleep(time.Second)

	if cs.cleanupCounter != currentCleanupCounter {
		return
	}

	maxSize := cs.maxSize
	currentSize := cs.store.Size()

	if currentSize > maxSize {
		diff := currentSize - maxSize

		var entries EntryTuples

		for _, e := range cs.store.Entries() {
			key, value := e.Get()
			entries = append(entries, NewTuple(key, value.evaluate()))
		}

		sort.Stable(entries)

		for i := 0; i < diff; i++ {
			if i < len(entries) {
				key, _ := entries[i].Get()

				cs.store.Delete(key)
			}
		}
	}
}

func (cs *CacheStore[T]) Set(key string, value T) {
	cs.store.Set(key, &CacheEntry[T]{
		value:       value,
		accessCount: 0,
		timestamp:   time.Now().UnixMilli(),
	})

	if cs.maxSize > 0 {
		go cs.removeEntriesIfMaxSizeExceeded()
	}

	if cs.entryLifetime > 0 {
		timeout(func() {
			cs.store.Delete(key)
		}, cs.entryLifetime)
	}
}

func (cs *CacheStore[T]) Flush() {
	cs.store.Clear()
}

func (cs *CacheStore[T]) Get(key string) (T, error) {
	entry, err := cs.store.Get(key)

	if err != nil {
		var t T
		return t, err
	}

	entry.accessCount++
	return entry.value, nil
}

func (cs *CacheStore[T]) Keys() []string {
	return cs.store.Keys()
}

func (cs *CacheStore[T]) Values() []T {
	e := cs.store.Values()

	var values []T

	for _, ce := range e {
		values = append(values, ce.value)
	}

	return values
}

func (cs *CacheStore[T]) Entries() []Tuple[string, T] {
	e := cs.store.Entries()

	var entries []Tuple[string, T]

	for _, ce := range e {
		key, entry := ce.Get()
		entries = append(entries, Tuple[string, T]{key, entry.value})
	}

	return entries
}
