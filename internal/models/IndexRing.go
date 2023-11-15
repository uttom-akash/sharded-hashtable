package models

import (
	"math/rand"
	"sync"

	"scale.kv.store/internal/core"
)

type ShardMetadata struct {
	ShardId uint32
}

func NewShardMetadata() *ShardMetadata {
	return &ShardMetadata{
		ShardId: 0, //Todo
	}
}

type ShardGroupMetadata struct {
	writeLock    sync.Mutex
	ShardGroupId uint32
	Shards       []*ShardMetadata //TODO shard metadata
}

func NewShardGroup() *ShardGroupMetadata {
	shards := make([]*ShardMetadata, 1)
	shards[0] = NewShardMetadata()

	return &ShardGroupMetadata{
		Shards:       shards,
		ShardGroupId: rand.Uint32(),
	}
}

const NUMBER_OF_SPOTS = 1000000

type IndexRing struct {
	Ring []*ShardGroupMetadata
	Size uint32
}

func NewIndexRing() *IndexRing {
	initialSpot := NewShardGroup()

	indexRing := make([]*ShardGroupMetadata, NUMBER_OF_SPOTS)

	for id := range indexRing {
		indexRing[id] = initialSpot
	}

	return &IndexRing{
		Ring: indexRing,
		Size: NUMBER_OF_SPOTS,
	}
}

func (indexRing *IndexRing) getSpotId(key byte) uint32 {
	hashedKey := core.Get32MurmurHash([]byte{key})

	spotId := hashedKey % indexRing.Size

	return spotId
}

func (indexRing *IndexRing) getNextSpot(spotId uint32) *ShardGroupMetadata {
	return indexRing.Ring[spotId]
}
