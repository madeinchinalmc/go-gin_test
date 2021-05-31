package concurrentMap

import "sync/atomic"

type BucketStatus uint8

const (
	//BUCKET_STATUS_NORMAL 散列桶正常
	BUCKET_STATUS_NORMAL BucketStatus = 0

	//BUCKET_STATUS_NORMAL 散列桶过轻
	BUCKET_STATUS_UNDERWEIGHT BucketStatus = 1

	//BUCKET_STATUS_OVERWEIGHT 散列桶过重
	BUCKET_STATUS_OVERWEIGHT BucketStatus = 2
)

// 针对键 - 元素对的再分布器
// 散列段内键- 元素对分布不均时重新分布
type PairRedistributor interface {
	// 根据键-元素对总数和散列桶总数计算并更新阈值
	UpdateThreshold(pairTotal uint64, bucketNumber int)

	//检查散列桶状态
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)

	//实施键- 元素对的再分布
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

//PairRedistributor 默认实现
type myPairRedistributor struct {
	//loadFactor 装载因子
	loadFactor float64
	//upperThreshold 散列桶重量的上阈限，散列桶尺寸增至此会触发再散列
	upperThreshold uint64
	//overweightBucketCount 过重散列桶计数
	overweightBucketCount uint64
	//emptyBucketCount 空散列桶计数
	emptyBucketCount uint64
}

var bucketCountTemplate = `Bucket count : 
	pairTotal: %d
	bucketNumber: %d
	average: %f
	upperThreshold: %d
	emptyBucketCount: %d
`

func (m *myPairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	if average <= 100 {
		average = 100
	}

	atomic.StoreUint64(&m.upperThreshold, uint64(average*m.loadFactor))
}

var bucketStatusTemplate = `Check bucket status: 
    pairTotal: %d
    bucketSize: %d
    upperThreshold: %d
    overweightBucketCount: %d
    emptyBucketCount: %d
    bucketStatus: %d
	
`

func (m *myPairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	if bucketSize > DEFAULT_BUCKET_MAX_SIZE ||
		bucketSize >= atomic.LoadUint64(&m.upperThreshold) {
		atomic.AddUint64(&m.overweightBucketCount, 1)
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	if bucketSize == 0 {
		atomic.AddUint64(&m.emptyBucketCount, 1)
	}
	return
}

var redistributionTemplate = `Redistributing: 
    bucketStatus: %d
    currentNumber: %d
    newNumber: %d
`

func (m *myPairRedistributor) Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber
	switch bucketStatus {
	case BUCKET_STATUS_OVERWEIGHT:
		if atomic.LoadUint64(&m.overweightBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber << 1
	case BUCKET_STATUS_UNDERWEIGHT:
		if currentNumber < 100 ||
			atomic.LoadUint64(&m.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber >> 1
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}
	if newNumber == currentNumber {
		atomic.StoreUint64(&m.overweightBucketCount, 0)
		atomic.StoreUint64(&m.emptyBucketCount, 0)
		return nil, false
	}
	var pairs []Pair
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}
	if newNumber > currentNumber {
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		buckets = make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}
	var count int
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
		count++
	}
	atomic.StoreUint64(&m.overweightBucketCount, 0)
	atomic.StoreUint64(&m.emptyBucketCount, 0)
	return buckets, true
}

// loadFactor 散列桶负载因子  bucketNumber 散列桶数量
func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &myPairRedistributor{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}
