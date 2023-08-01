package radix

import (
	"sync"
)

type safeTree struct {
	mutex *sync.RWMutex
	value *tree
}

func newSafeTree() Tree {
	return &safeTree{mutex: new(sync.RWMutex), value: newTree()}
}

func (st *safeTree) Lock() {
	st.mutex.Lock()
}

func (st *safeTree) Unlock() {
	st.mutex.Unlock()
}

func (st *safeTree) RLock() {
	st.mutex.RLock()
}

func (st *safeTree) RUnlock() {
	st.mutex.RUnlock()
}
func (st *safeTree) Insert(key Key, value Value) (Value, bool) {
	st.Lock()
	oldValue, updated := st.value.Insert(key, value)
	st.Unlock()
	return oldValue, updated
}

func (st *safeTree) Delete(key Key) (Value, bool) {
	st.Lock()
	value, del := st.value.Delete(key)
	st.Unlock()

	return value, del
}

func (st *safeTree) Search(key Key) (Value, bool) {
	st.RLock()
	val, ok := st.value.Search(key)
	st.RUnlock()
	return val, ok
}

func (st *safeTree) Minimum() (value Value, found bool) {
	return st.value.Minimum()
}

func (st *safeTree) Maximum() (value Value, found bool) {
	st.RLock()
	val, ok := st.value.Maximum()
	st.RUnlock()
	return val, ok
}

func (st *safeTree) Size() int {
	st.RLock()
	n := st.value.Size()
	st.RUnlock()
	return n
}

func (st *safeTree) ForEach(callback Callback, opts ...int) {
	st.RLock()
	st.value.ForEach(callback, opts...)
	st.RUnlock()
}

func (st *safeTree) ForEachPrefix(key Key, callback Callback) {
	st.RLock()
	st.value.ForEachPrefix(key, callback)
	st.RUnlock()
}

// Iterator pattern
func (st *safeTree) Iterator(opts ...int) Iterator {
	return st.value.Iterator(opts...)
}

func (st *safeTree) Clean() {
	st.Lock()
	defer st.Unlock()
	st.value.Clean()
}
