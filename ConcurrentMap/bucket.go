package concurrentMap

import (
	"sync"
	"sync/atomic"
)

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

// 并发安全的散列桶的实现类型
type bucket struct {
	// 键- 元素 对列表的表头
	firstValue atomic.Value
	size       uint64
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	if target != nil {
		target.SetElement(p.Element())
		return false, nil
	}
	p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}

func (b *bucket) Get(key string) Pair {
	panic("implement me")
}

func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	panic("implement me")
}

func (b *bucket) Clear(lock sync.Locker) {
	panic("implement me")
}

func (b *bucket) Size() uint64 {
	panic("implement me")
}

func (b *bucket) String() string {
	panic("implement me")
}

var placeholder Pair = &pair{}

// newBucket 会创建一个Bucket类型的实例。
func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}
