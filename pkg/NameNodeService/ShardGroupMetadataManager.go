package coordinationservice

import (
	"scale.kv.store/internal/metadata"
)

type ShardGroupMetadataManager struct {
	shardGroups []*metadata.ShardGroupMetadata
}

func NewShardGroupMetadataManager() *ShardGroupMetadataManager {

	shardGroups := make([]*metadata.ShardGroupMetadata, 1)
	shardGroups[0] = metadata.NewShardGroup(0)

	return &ShardGroupMetadataManager{
		shardGroups: shardGroups,
	}
}

func (shardGroupManager *ShardGroupMetadataManager) GetShardGroup(shardGroupId uint32) *metadata.ShardGroupMetadata {
	//shard group will always be available at index since shardGroupId has been fetched from filled
	return shardGroupManager.shardGroups[shardGroupId]
}
