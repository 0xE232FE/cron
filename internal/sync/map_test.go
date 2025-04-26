package sync

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestHas(t *testing.T) {
	m := Map[string, int]{}
	assert.False(t, m.Has("test"))
	m.Store("test", 1)
	assert.True(t, m.Has("test"))
	m.Delete("test")
	assert.False(t, m.Has("test"))
}

func TestRangeKeys(t *testing.T) {
	m := Map[string, int]{}
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	keys := make([]string, 0)
	m.RangeKeys(func(k string) bool {
		keys = append(keys, k)
		return true
	})
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestRangeValues(t *testing.T) {
	m := Map[string, int]{}
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	values := make([]int, 0)
	m.RangeValues(func(k int) bool {
		values = append(values, k)
		return true
	})
	sort.Ints(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}
