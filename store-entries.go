package stores

type storeEntries[T any] map[string]*T

func (entries *storeEntries[T]) find(key string) *T {
	return (*entries)[key]
}

func (entries *storeEntries[T]) set(key string, content *T) {
	(*entries)[key] = content
}

func (entries *storeEntries[T]) delete(key string) {
	delete(*entries, key)
}

func (entries *storeEntries[T]) clear() {
	*entries = make(storeEntries[T])
}
