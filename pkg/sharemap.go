package pkg

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const numPartitions = 16

type ShardedMap[K comparable, V any] struct {
	maps      [numPartitions]map[K]V
	partition [numPartitions]sync.RWMutex
}

func NewShardedMap[K comparable, V any]() *ShardedMap[K, V] {
	sm := &ShardedMap[K, V]{}
	for i := 0; i < numPartitions; i++ {
		sm.maps[i] = make(map[K]V)
	}
	return sm
}

func (sm *ShardedMap[K, V]) getPartition(key K) int {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%v", key)))
	return int(h.Sum64() % uint64(numPartitions))
}

func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
	partition := sm.getPartition(key)
	sm.partition[partition].RLock()
	defer sm.partition[partition].RUnlock()
	value, ok := sm.maps[partition][key]
	return value, ok
}

func (sm *ShardedMap[K, V]) Set(key K, value V) {
	partition := sm.getPartition(key)
	sm.partition[partition].Lock()
	defer sm.partition[partition].Unlock()
	sm.maps[partition][key] = value
}

func (sm *ShardedMap[K, V]) Delete(key K) {
	partition := sm.getPartition(key)
	sm.partition[partition].Lock()
	defer sm.partition[partition].Unlock()
	delete(sm.maps[partition], key)
}
