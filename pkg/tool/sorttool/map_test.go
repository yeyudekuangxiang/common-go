package sorttool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

func TestStringSortMap(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"c": "c",
		"b": "b",
		"1": "1",
		"A": "A",
		"C": "C",
		"B": "B",
	}

	wantKeys := []string{
		"1",
		"A",
		"B",
		"C",
		"a",
		"b",
		"c",
	}
	i := 0
	Map(m, func(key interface{}) {
		assert.Equal(t, wantKeys[i], key.(string))
		i++
	})
}
func TestIntSortMap(t *testing.T) {
	m := map[int]string{
		1:    "a",
		3:    "c",
		8:    "b",
		9:    "1",
		100:  "A",
		-1:   "C",
		0:    "B",
		-100: "B",
	}

	wantKeys := []int{
		-100,
		-1,
		0,
		1,
		3,
		8,
		9,
		100,
	}
	i := 0
	Map(m, func(key interface{}) {
		assert.Equal(t, wantKeys[i], key.(int))
		i++
	})
}
func TestPointKeyMap(t *testing.T) {
	m := make(map[*int]int)
	keys := make([]int, 0)
	for i := 0; i < 10; i++ {
		k := rand.Intn(100)
		m[&k] = k
		keys = append(keys, k)
	}
	sort.Ints(keys)
	i := 0
	Map(m, func(key interface{}) {
		assert.Equal(t, keys[i], *key.(*int))
		i++
	})
}
func TestPointMap(t *testing.T) {
	m := make(map[int]int)
	keys := make([]int, 0)
	for i := 0; i < 10; i++ {
		k := rand.Int()
		m[k] = k
		keys = append(keys, k)
	}
	sort.Ints(keys)
	i := 0
	Map(&m, func(key interface{}) {
		assert.Equal(t, keys[i], key)
		i++
	})
}
func TestStruct(t *testing.T) {
	m := struct {
		Name string
	}{
		Name: "test",
	}
	Map(&m, func(key interface{}) {
		fmt.Println(key)
	})
}
func TestFloatMap(t *testing.T) {
	m := make(map[int]int)
	keys := make([]int, 0)
	for i := 0; i < 10; i++ {
		k := rand.Int()
		m[k] = k
		keys = append(keys, k)
	}
	sort.Ints(keys)
	i := 0
	Map(&m, func(key interface{}) {
		assert.Equal(t, keys[i], key)
		i++
	})
}
