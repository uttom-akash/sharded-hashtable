package models

import (
	"fmt"
	"sync"
	"time"

	"scale.kv.store/internal/core"
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
		bloomFilter: core.NewBloomFilter(10000, 0.01, "optimal"),
		// bloom filter : https://hur.st/bloomfilter/?n=50000&p=0.00002&m=&k=
	}
}

func (shard *Shard) Get(key byte) *Value {

	if !shard.bloomFilter.Check([]byte{key}) {
		// key doesn't exist in this shard
		return nil
	}

	keyObject := NewKey(key)

	for _, bucket := range shard.Buckets {
		if bucket == nil {
			continue
		}

		startTime := time.Now()

		slot := bucket.Get(keyObject)

		if slot != nil {
			endTime := time.Now()

			duration := endTime.Sub(startTime)

			fmt.Printf("Get execution took %s\n", duration)

			return slot.value
		}
	}
	return nil
}

func (shard *Shard) Put(key byte, value byte) *Value {
	keyObject := NewKey(key)
	valueObject := NewValue(value)

	notExist := !shard.bloomFilter.Check([]byte{key})

	shard.writeLock.Lock()
	defer shard.writeLock.Unlock()

	for index, bucket := range shard.Buckets {
		if bucket == nil {
			bucket = NewBucket()
			shard.Buckets[index] = bucket
		}
		var slot *Slot

		if notExist {
			slot = bucket.PutNewKey(keyObject, valueObject)
		} else {
			slot = bucket.Put(keyObject, valueObject)
		}

		if slot != nil {
			shard.bloomFilter.Add([]byte{slot.key.Key})
			return slot.value
		}
	}
	return nil
}

func (shard *Shard) Delete(key byte) bool {
	notExist := !shard.bloomFilter.Check([]byte{key})

	if notExist {
		return false
	}

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
	return false
}
