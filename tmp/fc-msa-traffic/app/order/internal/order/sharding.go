package order

import (
	"hash/fnv"
)

// ShardRouter determines which shard to use based on a string key (e.g., user ID).
type ShardRouter struct {
	ShardCount int // Total number of shards
}

// NewShardRouter creates a new ShardRouter with the given shard count.
func NewShardRouter(shardCount int) *ShardRouter {
	if shardCount <= 0 {
		panic("shardCount must be greater than zero")
	}
	return &ShardRouter{ShardCount: shardCount}
}

// GetShard returns the index of the shard for the given string key.
// It uses FNV hashing to generate a consistent hash.
func (r *ShardRouter) GetShard(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	hashValue := h.Sum32()
	return int(hashValue) % r.ShardCount
}
