package models

import (
	"scale.kv.store/internal/core"
	"sync"
)

type Shard struct {
	Version     *Version
	bloomFilter *core.BloomFilter
	Buckets     [4096]*Bucket
	writeLock   sync.Mutex
}

func NewShard() *Shard {
	return &Shard{
		Version:     NewVersion(),
		bloomFilter: core.NewBloomFilter(50000, 0.00002, "optimal"),
	}
}

func (shard *Shard) Get(key byte) *Value {

	if !shard.bloomFilter.Check([]byte{key}) {
		return nil
	}

	keyObject := NewKey(key)
	for _, bucket := range shard.Buckets {
		if bucket == nil {
			continue
		}

		slot := bucket.Get(keyObject)
		if slot != nil {
			return slot.value
		}
	}
	return nil
}

func (shard *Shard) Put(key byte, value byte) *Value {
	keyObject := NewKey(key)
	valueObject := NewValue(value)

	shard.writeLock.Lock()
	defer shard.writeLock.Unlock()

	for index, bucket := range shard.Buckets {
		if bucket == nil {
			bucket = NewBucket()
			shard.Buckets[index] = bucket
		}

		slot := bucket.Put(keyObject, valueObject)
		if slot != nil {
			shard.bloomFilter.Add([]byte{slot.key.Key})
			return slot.value
		}
	}
	return nil
}

func (shard *Shard) Delete(key byte) bool {

	keyObject := NewKey(key)
	for _, bucket := range shard.Buckets {
		if bucket == nil {
			continue
		}

		deleted := bucket.Delete(keyObject)
		if deleted {
			return true
		}
	}
	return true
}
