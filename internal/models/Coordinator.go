package models

import (
	"sync"

	"scale.kv.store/internal/core"
)

const NUMBER_OF_SPOTS = 1000000

type Spot struct {
	writeLock sync.Mutex
	Shards    []*Shard
}

func NewSpot() *Spot {
	shards := make([]*Shard, 1)
	shards[0] = NewShard()

	return &Spot{
		Shards: shards,
	}
}

type Coordinator struct {
	IndexRing []*Spot
}

func NewCoordinator() *Coordinator {

	indexRing := make([]*Spot, NUMBER_OF_SPOTS)
	indexRing[0] = NewSpot()

	return &Coordinator{
		IndexRing: indexRing,
	}
}

func (coordinator *Coordinator) Get(key byte) *Result {

	spotId := coordinator.getSpotId(key)

	spot := coordinator.getNextSpot(spotId)

	var result *Result

	for _, shard := range spot.Shards {

		result := shard.Get(key)

		if result.Status == Found || result.Status == StopSearch {
			return result
		}
	}

	return result
}

func (coordinator *Coordinator) getNextSpot(spotId uint32) *Spot {

	iterator := NewIterator(int(spotId), NUMBER_OF_SPOTS)

	index := iterator.Current()

	for index != -1 {

		if coordinator.IndexRing[index] != nil {
			break
		}
		index = iterator.Next()
	}

	return coordinator.IndexRing[index]
}

func (coordinator *Coordinator) Put(key byte, value byte) *Result {

	spotId := coordinator.getSpotId(key)

	spot := coordinator.getNextSpot(spotId)

	spot.writeLock.Lock()
	defer spot.writeLock.Unlock()

	var bucketId *int
	var shardd *Shard
	var result *Result

	for _, shard := range spot.Shards {

		result := shard.SearchForWrite(key)

		if result.Status == Found || result.Status == StopSearch {

			result = shard.Put(key, value, *result.BucketId)

			return result

		} else {
			shardd = shard
			bucketId = result.BucketId
		}
	}

	if bucketId != nil && shardd != nil {
		result = shardd.Put(key, value, *bucketId)
	}

	return result
}

func (*Coordinator) getSpotId(key byte) uint32 {
	hashedKey := core.Get32MurmurHash([]byte{key})

	spotId := hashedKey % NUMBER_OF_SPOTS

	return spotId
}
