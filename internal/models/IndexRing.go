package models

import (
	"sync"

	"scale.kv.store/internal/core"
)

const NUMBER_OF_SPOTS = 1000000

type ShardGroupMetadata struct {
	writeLock    sync.Mutex
	ShardGroupId uint32
	Shards       []*Shard //TODO shard metadata
}

func NewShardGroup() *ShardGroupMetadata {
	shards := make([]*Shard, 1)
	shards[0] = NewShard()

	return &ShardGroupMetadata{
		Shards: shards,
	}
}

type IndexRing struct {
	IndexRing []*ShardGroupMetadata
}

func NewIndexRing() *IndexRing {
	initialSpot := NewShardGroup()

	indexRing := make([]*ShardGroupMetadata, NUMBER_OF_SPOTS)

	for id := range indexRing {
		indexRing[id] = initialSpot
	}

	return &IndexRing{
		IndexRing: indexRing,
	}
}

func (*IndexRing) getSpotId(key byte) uint32 {
	hashedKey := core.Get32MurmurHash([]byte{key})

	spotId := hashedKey % NUMBER_OF_SPOTS

	return spotId
}

func (indexRing *IndexRing) getNextSpot(spotId uint32) *ShardGroupMetadata {
	return indexRing.IndexRing[spotId]
}
