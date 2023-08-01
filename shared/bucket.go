package shared

import "time"

type ShareBucket struct {
	name   string
	bucket Cache
}

func (shm *ShareBucket) Get(key string) (interface{}, error) {
	return shm.bucket.Get(key)
}

func (shm *ShareBucket) Del(key string) {
	shm.bucket.Remove(key)
}

func (shm *ShareBucket) Set(key string, val interface{}, ttl int) error {
	if ttl <= 0 {
		return shm.bucket.Set(key, val)
	}

	expire := time.Duration(ttl) * time.Second
	return shm.bucket.SetWithExpire(key, val, expire)
}

func (shm *ShareBucket) Clear() {

}

func (shm *ShareBucket) reg() {

}

func newShareBucket(name string, bkt Cache) *ShareBucket {
	s := &ShareBucket{
		name:   name,
		bucket: bkt,
	}
	return s
}
