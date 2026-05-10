package hash

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
	"github.com/aayush0325/consistent-hashing/internal/services/metrics"
	"github.com/spaolacci/murmur3"
)

func Set(key string, val string, ttl uint64, ctx context.Context) error {
	if len(coordinator.Ring) == 0 {
		return fmt.Errorf("ring is empty")
	}

	hash := murmur3.Sum64([]byte(key))

	index := upperBound(hash)

	if index == len(coordinator.Ring) {
		index = 0
	}

	maxTries := config.App.MaxTries
	if maxTries <= 0 {
		maxTries = len(coordinator.Ring)
	}

	tries := 0

	for tries < maxTries {
		node := coordinator.Ring[index].ParentNode

		if coordinator.GlobalState[node].IsHealthy {
			client := coordinator.GlobalState[node].Client

			log.Printf("Attempting to write key %s to node %s (try %d/%d)", key, node, tries+1, maxTries)
			_, err := client.Set(ctx, key, val, time.Duration(ttl)*time.Second).Result()
			if err == nil {
				log.Printf("Successfully wrote key %s to node %s", key, node)
				if config.App.ReplicationFactor > 1 {
					log.Printf("Starting background replication for key %s", key)
					go replicateAsync(key, val, ttl, (index+1)%len(coordinator.Ring), maxTries-tries-1)
				}
				metrics.M.IncCacheHit("set")
				return nil
			}
			metrics.M.IncCacheMiss("set")
			log.Printf("Failed to write key %s to node %s: %v", key, node, err)
		} else {
			metrics.M.IncCacheMiss("set")
			log.Printf("Node %s is unhealthy, skipping", node)
		}

		index = (index + 1) % len(coordinator.Ring)
		tries++
	}

	return fmt.Errorf("no healthy nodes available")
}

func replicateAsync(key string, val string, ttl uint64, startIndex int, maxTries int) {
	ctx := context.Background()
	tries := 0
	written := 0
	index := startIndex

	for tries < maxTries && written < config.App.ReplicationFactor-1 {
		node := coordinator.Ring[index].ParentNode

		if coordinator.GlobalState[node].IsHealthy {
			client := coordinator.GlobalState[node].Client

			log.Printf("Background: Attempting to replicate key %s to node %s", key, node)
			_, err := client.Set(ctx, key, val, time.Duration(ttl)*time.Second).Result()
			if err == nil {
				log.Printf("Background: Successfully replicated key %s to node %s", key, node)
				written++
			} else {
				log.Printf("Background: Failed to replicate key %s to node %s: %v", key, node, err)
			}
		}

		index = (index + 1) % len(coordinator.Ring)
		tries++
	}
	log.Printf("Background replication for key %s finished. Written to %d additional nodes", key, written)
}
