package slicemap

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestBasicPutGet(t *testing.T) {

	sm := SliceMap[string, int]{}
	if sm.Len() != 0 {
		t.Errorf("Expected length 0")
	}
	sm.Put([]string{}, 1)
	if sm.Len() != 1 {
		t.Errorf("Expected length 1")
	}
	v, ok := sm.Get([]string{})
	if !ok {
		t.Errorf("Expected ok=true")
	}
	if v != 1 {
		t.Errorf("Expecting 1 got %d", v)
	}
	sm.Put([]string{"a", "b", "c"}, 2)
	sm.Put([]string{"a", "b", "d"}, 3)
	sm.Put([]string{"a", "b"}, 4)
	if v, ok := sm.Get([]string{"a", "b"}); !ok || v != 4 {
		t.Errorf("Wrong a,b: %v", v)
	}
	if v, ok := sm.Get([]string{"a", "b", "c"}); !ok || v != 2 {
		t.Errorf("Wrong a,b,c: %v", v)
	}
	if v, ok := sm.Get([]string{"a", "b", "d"}); !ok || v != 3 {
		t.Errorf("Wrong a,b,d: %v", v)
	}
}

func TestRandom(t *testing.T) {
	sm := SliceMap[string, int]{}
	keys := make([][]string, 1000)
	for i := 0; i < len(keys); i++ {
		for {
			key := make([]string, 0, 5)
			n := rand.Intn(5)
			for x := 0; x < n; x++ {
				key = append(key, strconv.Itoa(rand.Int()))
			}
			found := false
			for ki := 0; ki < i; ki++ {
				if reflect.DeepEqual(keys[ki], key) {
					found = true
					break
				}
			}
			if !found {
				keys[i] = key
				break
			}
		}
	}
	for i, k := range keys {
		sm.Put(k, i)
	}
	if sm.Len() != len(keys) {
		t.Errorf("Wrong len: %d", sm.Len())
	}
	for i := 0; i < len(keys); i++ {
		if v, _ := sm.Get(keys[i]); v != i {
			t.Errorf("Wrong value for %v %v", keys[i], v)
		}
	}
	foundValue := make(map[int]struct{})
	sm.ForEach(func(k []string, v int) bool {
		found := false
		for i := 0; i < len(keys); i++ {
			if reflect.DeepEqual(keys[i], k) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Key %v not found", k)
		}
		foundValue[v] = struct{}{}
		return true
	})
	if len(foundValue) != len(keys) {
		t.Errorf("Found: %d", len(foundValue))
	}

	for len(keys) > 0 {
		n := rand.Intn(len(keys))
		sm.Delete(keys[n])
		_, exists := sm.Get(keys[n])
		if exists {
			t.Errorf("Still exists: %v", keys[n])
		}
		keys = append(keys[:n], keys[n+1:]...)
		if len(keys) != sm.Len() {
			t.Errorf("Expecting %d got %d", len(keys), sm.Len())
		}
	}
}
