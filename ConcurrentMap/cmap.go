package concurrentMap

type ConcurrentMap interface {
	Concurrency() int
	Put(key string,element interface{}) (bool,error)
	Get(key string) interface{}
	Delete(key string) bool
	Len() uint64
}
