package hash

import (
	"context"
	"fmt"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
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

			res, err := client.Get(ctx, string(s)).Result()
			if err != nil {
				return nil, err
			}

			return []byte(res), nil
		}

		index = (index + 1) % len(coordinator.Ring)
		tries++
	}

	return nil, fmt.Errorf("no healthy nodes available")
}
