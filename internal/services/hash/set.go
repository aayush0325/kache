package hash

import (
	"context"
	"fmt"
	"time"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
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

			_, err := client.Set(ctx, key, val, time.Duration(ttl)).Result()
			if err != nil {
				return err
			}

			return nil
		}

		index = (index + 1) % len(coordinator.Ring)
		tries++
	}

	return fmt.Errorf("no healthy nodes available")
}
