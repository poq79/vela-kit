package dict

import "sync/atomic"

type MemDict struct {
	size   int32
	offset int32
	mem    *Memory
}

func (m *MemDict) Reset() error {
	atomic.StoreInt32(&m.offset, 0)
	return nil
}

func (m *MemDict) Close() error {
	return nil
}

func (m *MemDict) Next() bool {
	return atomic.AddInt32(&m.offset, 1) <= m.size
}

func (m *MemDict) Text() string {
	if m.size == 0 {
		return ""
	}

	if m.offset == 0 {
		return m.mem.value[m.offset]
	}

	if m.offset > m.size {
		return ""
	}

	return m.mem.value[m.offset-1]
}

func (m *MemDict) Done() {
	atomic.StoreInt32(&m.offset, m.size)
}
