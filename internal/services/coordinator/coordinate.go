package coordinator

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aayush0325/consistent-hashing/internal/config"
	"github.com/redis/go-redis/v9"
)

type NodeState struct {
	IsHealthy bool          `json:"isHealthy"`
	Client    *redis.Client `json:"-"`
}

type State map[string]*NodeState

var GlobalState = make(State)

func dialAddr(addr string) error {
	dialer := &net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func PingAllNodes(ctx context.Context) {
	log.Printf("Pinging all nodes")
	for addr, state := range GlobalState {
		if state.Client == nil {
			// if we cannot establish a tcp connection then the node must be
			// unhealthy
			if err := dialAddr(addr); err != nil {
				if state.IsHealthy {
					log.Printf("node %s became unhealthy: %v", addr, err)
					GlobalState[addr].IsHealthy = false
				}
				continue
			}

			// tcp can be established, create a redis connection
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: "",
				DB:       0,
			})
			GlobalState[addr].Client = client
			if !state.IsHealthy {
				log.Printf("node %s became healthy again", addr)
			}
			GlobalState[addr].IsHealthy = true
			continue
		}

		// tcp was previously there now broken, hence the node is unhealthy
		if err := dialAddr(addr); err != nil {
			state.Client.Close()
			GlobalState[addr].Client = nil
			if state.IsHealthy {
				log.Printf("node %s became unhealthy: %v", addr, err)
			}
			GlobalState[addr].IsHealthy = false
			continue
		}

		// tcp is there but the ping failed
		err := state.Client.Ping(ctx).Err()
		if err != nil {
			state.Client.Close()
			GlobalState[addr].Client = nil
			if state.IsHealthy {
				log.Printf("node %s became unhealthy: %v", addr, err)
			}
			GlobalState[addr].IsHealthy = false
		} else {
			// otherwise node is healthy
			if !state.IsHealthy {
				log.Printf("node %s is healthy", addr)
			}
			GlobalState[addr].IsHealthy = true
		}
	}
}

func readRedisNodes() {
	redisUrls := os.Getenv("REDIS_URLS")
	redisUrlsArr := strings.Split(redisUrls, ",")
	for _, value := range redisUrlsArr {
		addr := strings.TrimSpace(value)
		GlobalState[addr] = &NodeState{
			IsHealthy: false,
			Client:    nil,
		}

		log.Printf("registered node %s (health pending)", addr)
	}
}

func Background(ctx context.Context) {
	readRedisNodes()
	createRing()

	ticker := time.NewTicker(config.App.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			PingAllNodes(ctx)
		case <-ctx.Done():
			return
		}
	}
}
