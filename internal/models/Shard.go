package models

import (
	"scale.kv.store/internal/core"
)

const NUMBER_OF_BUCKET = 4096

type Shard struct {
	Version     *Version
	bloomFilter *core.BloomFilter
	Buckets     [NUMBER_OF_BUCKET]*Bucket
}

func NewShard() *Shard {
	return &Shard{
		Version: NewVersion(),
		//TODO: reduce the size
		//TODO: check alternative like cuckoo filter
		bloomFilter: core.NewBloomFilter(10000, 0.01, "optimal"),
		// bloom filter : https://hur.st/bloomfilter/?n=50000&p=0.00002&m=&k=
	}
}

func (shard *Shard) Get(key byte) *Result {

	//quick check
	if shard.bloomFilter.DoesNotExist([]byte{key}) {

		return NewContinueSearchResult()
	}

	keyObject := NewKey(key)

	hashedBucketId := core.Get16MurmurHash([]byte{keyObject.Key})

	iterator := NewIterator(int(hashedBucketId), NUMBER_OF_BUCKET)

	bucketId := iterator.Current()
	for bucketId != -1 {

		bucket := shard.Buckets[bucketId]

		if bucket == nil {
			// path end
			return NewStopSearchResult()
		}

		slot := bucket.Get(keyObject)

		if slot != nil {
			return NewFoundResult(*slot.value)
		}

		if bucket.HasEmptySlots() {
			// path end
			return NewStopSearchResult()
		}

		bucketId = iterator.Next()
	}

	//Couldn't found and also buckets are full
	return NewContinueSearchResult()
}

func (shard *Shard) SearchForWrite(key byte) *Result {

	//quick check
	//TODO: keep track of empty and deleted bucket
	// if shard.bloomFilter.DoesNotExist([]byte{key}) {

	// 	return NewContinueSearchResult()
	// }

	keyObject := NewKey(key)

	hashedBucketId := core.Get16MurmurHash([]byte{keyObject.Key})

	iterator := NewIterator(int(hashedBucketId), NUMBER_OF_BUCKET)

	bucketId := iterator.Current()

	var bucketHavingDeletedSlot int

	for bucketId != -1 {

		bucket := shard.Buckets[bucketId]

		if bucket == nil {
			// path end
			return NewStopSearchResultWithBucket(bucketId)
		}

		slot := bucket.Get(keyObject)

		if slot != nil {
			return NewFoundResultWithBucket(*slot.value, bucketId)
		}

		if bucket.HasEmptySlots() {
			// path end
			return NewStopSearchResultWithBucket(bucketId)
		}

		if bucket.HasDeletedSlots() { //TODO: set first deleted slot
			bucketHavingDeletedSlot = bucketId
		}

		bucketId = iterator.Next()
	}

	//Couldn't found and also buckets are full
	return NewContinueSearchResultWithBucket(bucketHavingDeletedSlot)
}

func (shard *Shard) Put(key byte, value byte, bucketId int) *Result {

	keyObject := NewKey(key)

	valueObject := NewValue(value)

	doesNotExist := shard.bloomFilter.DoesNotExist([]byte{key})

	if shard.Buckets[bucketId] == nil {
		shard.Buckets[bucketId] = NewBucket()
	}
	bucket := shard.Buckets[bucketId]

	var slot *Slot

	if doesNotExist {
		slot = bucket.PutNewKey(keyObject, valueObject)
	} else {
		slot = bucket.Put(keyObject, valueObject)
	}

	if slot != nil {
		shard.bloomFilter.Add([]byte{slot.key.Key})

		return NewAddedOrUpdatedResult(*slot.value)
	}

	return NewContinueSearchResult()
}

func (shard *Shard) Delete(key byte) *Result {
	//quick check
	if shard.bloomFilter.DoesNotExist([]byte{key}) {

		return NewContinueSearchResult()
	}

	keyObject := NewKey(key)

	hashedBucketId := core.Get16MurmurHash([]byte{keyObject.Key})

	iterator := NewIterator(int(hashedBucketId), NUMBER_OF_BUCKET)

	bucketId := iterator.Current()
	for bucketId != -1 {

		bucket := shard.Buckets[bucketId]

		if bucket == nil {
			// path end
			return NewStopSearchResult()
		}

		deleted := bucket.Delete(keyObject)
		if deleted {
			return NewDeletedResult()
		}

		if bucket.HasEmptySlots() {
			// path end
			return NewStopSearchResult()
		}

		bucketId = iterator.Next()
	}

	return NewContinueSearchResult()
}
