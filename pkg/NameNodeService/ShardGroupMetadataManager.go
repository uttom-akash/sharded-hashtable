package coordinationservice

import (
	"scale.kv.store/internal/metadata"
)

type ShardGroupMetadataManager struct {
	shardGroups []*metadata.ShardGroupMetadata
}

func NewShardGroupMetadataManager() *ShardGroupMetadataManager {

	shards := make([]*metadata.ShardGroupMetadata, 1)
	shards[0] = metadata.NewShardGroup(0)

	return &ShardGroupMetadataManager{
		shardGroups: shards,
	}
}

func (shardManager *ShardGroupMetadataManager) GetShardGroup(shardGroupId uint32) *metadata.ShardGroupMetadata {
	//check index
	return shardManager.shardGroups[shardGroupId]
}
