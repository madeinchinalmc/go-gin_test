package concurrentMap

type BucketStatus struct {
}

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
