package coordinator

import (
	"fmt"
	"log"
	"sort"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/spaolacci/murmur3"
)

type VirtualNode struct {
	ParentNode        string
	virtualNodeNumber uint64
	Hash              uint64
}

var Ring = make([]VirtualNode, 0)

func createRing() {
	Ring = Ring[:0]

	for server := range GlobalState {
		for vnodeNumber := range config.App.VNodesPerNode {
			vnodeId := fmt.Sprintf("%s_#%d", server, vnodeNumber)

			hash := murmur3.Sum64([]byte(vnodeId))

			Ring = append(Ring, VirtualNode{
				ParentNode:        server,
				virtualNodeNumber: uint64(vnodeNumber),
				Hash:              hash,
			})
		}
	}

	sort.Slice(Ring, func(i, j int) bool {
		return Ring[i].Hash < Ring[j].Hash
	})

	log.Printf("Hash ring created with %d virtual nodes", len(Ring))
}
