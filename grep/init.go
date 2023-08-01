package grep

var pool *Cache

func init() {
	pool = newLRU(4096)
}
