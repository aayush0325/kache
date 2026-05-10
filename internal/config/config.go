package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	VNodesPerNode     int
	PingInterval      time.Duration
	MaxTries          int
	ReplicationFactor int
}

var App = load()

func load() Config {
	return Config{
		VNodesPerNode:     getEnvInt("VNODES_PER_NODE", 100),
		PingInterval:      time.Duration(getEnvInt("PING_INTERVAL_SECONDS", 3)) * time.Second,
		MaxTries:          getEnvInt("MAX_TRIES", 0),
		ReplicationFactor: getEnvInt("REPLICATION_FACTOR", 2),
	}
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
