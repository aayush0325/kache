package hash

import (
	"context"
	"fmt"
	"log"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
	"github.com/aayush0325/consistent-hashing/internal/services/metrics"
	"github.com/spaolacci/murmur3"
)

func Get(s []byte, ctx context.Context) ([]byte, error) {
	if len(coordinator.Ring) == 0 {
		return nil, fmt.Errorf("ring is empty")
	}

	hash := murmur3.Sum64(s)

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

			log.Printf("Attempting to fetch key %s from node %s (try %d/%d)", string(s), node, tries+1, maxTries)
			res, err := client.Get(ctx, string(s)).Result()
			if err != nil {
				metrics.M.IncCacheMiss("get")
				log.Printf("Failed to fetch key %s from node %s: %v", string(s), node, err)
				index = (index + 1) % len(coordinator.Ring)
				tries++
				continue
			}

			log.Printf("Successfully fetched key %s from node %s", string(s), node)
			metrics.M.IncCacheHit("get")
			return []byte(res), nil
		} else {
			metrics.M.IncCacheMiss("get")
			log.Printf("Node %s is unhealthy, skipping", node)
		}

		index = (index + 1) % len(coordinator.Ring)
		tries++
	}

	return nil, fmt.Errorf("no healthy nodes available")
}
