package concurrentMap

import "sync/atomic"

type ConcurrentMap interface {
	Concurrency() int
	Put(key string, element interface{}) (bool, error)
	Get(key string) interface{}
	Delete(key string) bool
	Len() uint64
}

type Segment interface {
}

type myConcurrentMap struct {
	concurrency int
	segments    []Segment
	total       uint64
}

type PairRedistributor interface {
}

func NewConcurrentMap(concurrency int, pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is too large")
	}
	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = nil //todo
	}
	return cmap, nil
}

func (cmap *myConcurrentMap) Concurrency() int {
	return cmap.concurrency
}

func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	//todo
	return true, nil
}

func (cmap *myConcurrentMap) Get(key string) interface{} {
	//todo
	return nil
}

func (cmap *myConcurrentMap) Delete(key string) bool {
	//todo
	return false
}
func (cmap *myConcurrentMap) Len() uint64 {
	return atomic.LoadUint64(&cmap.total)
}
