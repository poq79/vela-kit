package lua

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type ExDataKV struct {
	key   string
	value interface{}
}

type ExData []ExDataKV

func (ed *ExData) Len() int {
	return len(*ed)
}

func (ed *ExData) Swap(i, j int) {
	a := *ed
	a[i], a[j] = a[j], a[i]
}

func (ed *ExData) Less(i, j int) bool {
	a := *ed
	return a[i].key < a[j].key
}

func (ed *ExData) Set(key string, value interface{}) {
	args := *ed

	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if key == kv.key {
			kv.value = value
			return
		}
		if kv.key == "" {
			kv.value = value
			return
		}
	}

	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		kv.key = key

		kv.value = value
		*ed = args
		return
	}

	kv := ExDataKV{}
	kv.key = key
	kv.value = value
	*ed = append(args, kv)

	//排序
	sort.Sort(ed)
}

func (ed *ExData) Get(key string) interface{} {

	a := *ed
	i, j := 0, ed.Len()
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		switch strings.Compare(key, a[h].key) {
		case 0:
			return a[h]
		case 1:
			i = h + 1
		case -1:
			j = h
		}
	}

	return nil
}

func (ed *ExData) Del(key string) {
	a := *ed
	n := len(a)
	for i := 0; i < n; i++ {
		kv := &a[i]
		if kv.key == key {
			a = append(a[:i], a[:i+1]...)
			goto DONE
		}
	}

DONE:
	*ed = a
}

func (ed *ExData) Reset() {
	*ed = (*ed)[:0]
}

type exUserKV struct {
	key string
	val LValue
}

type UserKV interface {
	LValue
	Get(string) LValue
	Set(string, LValue)
	V(string) (LValue, bool)
}

type userKV struct {
	data []exUserKV
}

func NewUserKV() UserKV {
	return &userKV{}
}

func (ukv *userKV) Len() int {
	return len(ukv.data)
}

func (ukv *userKV) cap() int {
	return cap(ukv.data)
}

func (ukv *userKV) Get(key string) LValue {
	n := ukv.Len()
	for i := 0; i < n; i++ {
		kv := &ukv.data[i]
		if kv.key == key {
			return kv.val
		}
	}
	return LNil
}

func (ukv *userKV) Set(key string, val LValue) {
	n := ukv.Len()
	for i := 0; i < n; i++ {
		kv := &ukv.data[i]
		if key == kv.key {
			kv.val = val
			return
		}
	}

	c := ukv.cap()
	if c > n {
		ukv.data = ukv.data[:n+1]
		kv := &ukv.data[n]
		kv.key = key
		kv.val = val
		return
	}

	kv := exUserKV{}
	kv.key = key
	kv.val = val

	ukv.data = append(ukv.data, kv)
}

func (ukv *userKV) V(key string) (LValue, bool) {
	n := ukv.Len()
	for i := 0; i < n; i++ {
		kv := &ukv.data[i]
		if kv.key == key {
			return kv.val, true
		}
	}
	return nil, false
}

func (ukv *userKV) String() string                     { return fmt.Sprintf("function: %p", ukv) }
func (ukv *userKV) Type() LValueType                   { return LTKv }
func (ukv *userKV) AssertFloat64() (float64, bool)     { return 0, false }
func (ukv *userKV) AssertString() (string, bool)       { return "", false }
func (ukv *userKV) AssertFunction() (*LFunction, bool) { return nil, false }
func (ukv *userKV) Peek() LValue                       { return ukv }

type safeUserKV struct {
	sync.RWMutex

	hook func(string) LValue
	data []exUserKV
}

func NewSafeUserKV() UserKV {
	return &safeUserKV{}
}

func (sukv *safeUserKV) Len() int {
	return len(sukv.data)
}

func (sukv *safeUserKV) cap() int {
	return cap(sukv.data)
}

func (sukv *safeUserKV) Swap(i, j int) {
	sukv.data[i], sukv.data[j] = sukv.data[j], sukv.data[i]
}

func (sukv *safeUserKV) Less(i, j int) bool {
	return sukv.data[i].key < sukv.data[j].key
}

func (sukv *safeUserKV) reset() {
	sukv.Lock()

	n := sukv.Len()
	for i := 0; i < n; i++ {
		sukv.data = nil
	}
	sukv.data = sukv.data[:0]
	sukv.Unlock()
}

func (sukv *safeUserKV) Set(key string, val LValue) {
	sukv.Lock()

	n := sukv.Len()
	c := sukv.cap()

	var newKV exUserKV
	for i := 0; i < n; i++ {
		kv := &sukv.data[i]
		if key == kv.key {
			kv.val = val
			goto done
		}
		if kv.key == "" {
			kv.val = val
			goto done
		}
	}

	if c > n {
		sukv.data = sukv.data[:n+1]
		kv := &sukv.data[n]
		kv.key = key
		kv.val = val
	}

	newKV = exUserKV{}
	newKV.key = key
	newKV.val = val
	sukv.data = append(sukv.data, newKV)

done:
	//排序
	sort.Sort(sukv)
	sukv.Unlock()
}

// 获取
func (sukv *safeUserKV) Get(key string) LValue {
	sukv.RLock()
	i, j := 0, sukv.Len()
	val := LNil
	for i < j {
		h := int(uint(i+j) >> 1)
		switch strings.Compare(key, sukv.data[h].key) {
		case 0:
			val = sukv.data[h].val
			goto done
		case 1:
			i = h + 1
		case -1:
			j = h
		}
	}

done:
	sukv.RUnlock()
	return val
}

func (sukv *safeUserKV) V(key string) (LValue, bool) {
	sukv.RLock()
	defer sukv.RUnlock()
	i, j := 0, sukv.Len()
	for i < j {
		h := int(uint(i+j) >> 1)
		switch strings.Compare(key, sukv.data[h].key) {
		case 0:
			return sukv.data[h].val, true

		case 1:
			i = h + 1
		case -1:
			j = h
		}
	}
	return nil, false
}

func (sukv *safeUserKV) String() string                     { return fmt.Sprintf("function: %p", sukv) }
func (sukv *safeUserKV) Type() LValueType                   { return LTSkv }
func (sukv *safeUserKV) AssertFloat64() (float64, bool)     { return 0, false }
func (sukv *safeUserKV) AssertString() (string, bool)       { return "", false }
func (sukv *safeUserKV) AssertFunction() (*LFunction, bool) { return nil, false }
func (sukv *safeUserKV) Peek() LValue                       { return sukv }
