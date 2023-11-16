package coordinationservice

import "scale.kv.store/internal/storageengine"

type ShardManager struct {
	Shards []*storageengine.Shard
}

func NewShardManager() *ShardManager {

	shards := make([]*storageengine.Shard, 1)
	shards[0] = storageengine.NewShard()

	return &ShardManager{
		Shards: shards,
	}
}

func (shardManager *ShardManager) GetShard(shardId uint32) *storageengine.Shard {
	//check index
	return shardManager.Shards[shardId]
}
