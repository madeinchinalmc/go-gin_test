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

func newSegment(bucketNumber int, pairRedistributor PairRedistributor) Segment {
	if bucketNumber < 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistributor == nil {
		pairRedistributor = nil //todo default pair redistributor
	}
	buckets := make([]Bucket, bucketNumber)
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = nil //todo bucket entity
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistributor,
	}
}

func (*segment) Put(p Pair) (bool, error) {
	return true, nil
}

func (*segment) Get(key string) Pair {
	return nil
}

func (*segment) GetWithHash(key string, keyHash uint64) Pair {
	return nil
}

func (*segment) Delete(key string) bool {
	return true
}
func (*segment) Size() uint64 {
	return 0
}
