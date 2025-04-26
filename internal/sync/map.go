package sync

import (
	"iter"
	"slices"
	"sync"
)

type Map[K comparable, V any] struct {
	m sync.Map
}

// Has returns either or not the map has the provided key
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.Load(key)
	return ok
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	var zeroV V
	if v, ok := m.m.Load(key); ok {
		return v.(V), ok
	}
	return zeroV, false
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	v, loaded := m.m.LoadOrStore(key, value)
	if loaded {
		return v.(V), loaded
	}
	return v.(V), false
}

func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	var zeroV V
	if v, loaded := m.m.LoadAndDelete(key); loaded {
		return v.(V), loaded
	}
	return zeroV, false
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Swap(key K, value V) (V, bool) {
	var zeroV V
	if previous, loaded := m.m.Swap(key, value); loaded {
		return previous.(V), loaded
	}
	return zeroV, false
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	// wrapper takes a "typed" callback and return a "untyped" callback
	// suitable for the underlying sync.Map Range function.
	wrapper := func(typedFn func(k K, v V) bool) func(any, any) bool {
		return func(k, v any) bool {
			typedKey, ok1 := k.(K)
			if !ok1 {
				panic("invalid key type")
			}
			typedValue, ok2 := v.(V)
			if !ok2 {
				panic("invalid value type")
			}
			return typedFn(typedKey, typedValue)
		}
	}
	m.m.Range(wrapper(f))
}

func (m *Map[K, V]) RangeKeys(f func(key K) bool) {
	m.Range(func(key K, _ V) bool {
		return f(key)
	})
}

func (m *Map[K, V]) RangeValues(f func(value V) bool) {
	m.Range(func(_ K, value V) bool {
		return f(value)
	})
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(yield)
	}
}

func (m *Map[K, V]) IterKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.RangeKeys(yield)
	}
}

func (m *Map[K, V]) IterValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.RangeValues(yield)
	}
}

func (m *Map[K, V]) Keys() (out []K) {
	return slices.Collect(m.IterKeys())
}

func (m *Map[K, V]) Values() (out []V) {
	return slices.Collect(m.IterValues())
}

func (m *Map[K, V]) Clear() {
	m.m.Clear()
}
