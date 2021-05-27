package concurrentMap

import "sync"

type Segment interface {
	Put(p Pair) (bool, error)
	Get(key string) Pair
	// 根据给定参数返回对应键 - 元素对
	GetWithHash(key string, keyHash uint64) Pair
	Delete(key string) bool
	Size() uint64
}

type segment struct {
	buckets           []Bucket
	bucketsLen        int
	pairTotal         uint64
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}
