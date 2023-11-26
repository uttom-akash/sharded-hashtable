package metadata

type ShardMetadata struct {
	ShardId uint32
}

func NewShardMetadata(shardId uint32) *ShardMetadata {
	return &ShardMetadata{
		ShardId: shardId, //Todo
	}
}
