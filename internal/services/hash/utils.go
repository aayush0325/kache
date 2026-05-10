package hash

import (
	"sort"

	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
)

func upperBound(v uint64) int {
	index := sort.Search(len(coordinator.Ring), func(i int) bool {
		return coordinator.Ring[i].Hash >= v
	})
	return index
}
