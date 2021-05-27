package concurrentMap

import "sync"

// 并发安全的散列桶接口
type Bucket interface {
	// put放入一个键 - 元素元素，调用此方法前lock了这里就不要把lock传入
	Put(p Pair, lock sync.Locker) (bool, error)

	// 获取指定 键 - 元素 对
	Get(key string) Pair

	// 返回第一个键 - 元素对
	GetFirstPair() Pair

	// 删除指定的 键 - 元素 对
	Delete(key string, lock sync.Locker) bool

	//清空散列桶
	Clear(lock sync.Locker)

	// 返回散列桶大小
	Size() uint64

	// 返回当前散列桶字符串表示形式
	String() string
}
