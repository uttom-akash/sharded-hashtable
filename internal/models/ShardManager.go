package models

type ShardManager struct {
	Shards []*Shard
}

func NewShardManager() *ShardManager {

	shards := make([]*Shard, 1)
	shards[0] = NewShard()

	return &ShardManager{
		Shards: shards,
	}
}

func (shardManager *ShardManager) GetShard(shardId uint32) *Shard {
	//check index
	return shardManager.Shards[shardId]
}
