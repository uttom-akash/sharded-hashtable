package metadata

import (
	"sync"
)

const INITIAL_SHARD_ID = 0

type ShardGroupMetadata struct {
	WriteLock    sync.Mutex
	ShardGroupId uint32
	Shards       []*ShardMetadata //TODO shard metadata
}

func NewShardGroup(shardGroupId uint32) *ShardGroupMetadata {
	shards := make([]*ShardMetadata, 1)
	shards[INITIAL_SHARD_ID] = NewShardMetadata(INITIAL_SHARD_ID)

	return &ShardGroupMetadata{
		Shards:       shards,
		ShardGroupId: shardGroupId,
	}
}
