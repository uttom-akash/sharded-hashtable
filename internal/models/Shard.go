package models

import (
	"sync"

	"scale.kv.store/internal/core"
)

const NUMBER_OF_BUCKET = 4096

type Shard struct {
	Version     *Version
	bloomFilter *core.BloomFilter
	Buckets     [NUMBER_OF_BUCKET]*Bucket
	writeLock   sync.Mutex
}

func NewShard() *Shard {
	return &Shard{
		Version:     NewVersion(),
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

	bucketId := iterator.Next()
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
	}

	//Couldn't found and also buckets are full
	return NewContinueSearchResult()
}

func (shard *Shard) Put(key byte, value byte) *Result {

	keyObject := NewKey(key)
	valueObject := NewValue(value)

	doesNotExist := shard.bloomFilter.DoesNotExist([]byte{key})

	shard.writeLock.Lock()
	defer shard.writeLock.Unlock()

	hashedBucketId := core.Get16MurmurHash([]byte{keyObject.Key})

	iterator := NewIterator(int(hashedBucketId), NUMBER_OF_BUCKET)

	bucketId := iterator.Next()
	for bucketId != -1 {

		bucket := shard.Buckets[bucketId]

		if bucket == nil {
			bucket = NewBucket()
			shard.Buckets[bucketId] = bucket
		}

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

		bucketId = iterator.Next()
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

	bucketId := iterator.Next()
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
	}

	return NewContinueSearchResult()
}
