package redis

import (
	"fmt"
	"log"

	"gopkg.in/redis.v3"
)

type RedisHandler struct {
	Client    *redis.Client
	Cluster   *redis.ClusterClient
	IsCluster bool
}

// Initialize Redis Handler (with optional password)
func NewRedisHandler(hosts []string, password string) *RedisHandler {
	if len(hosts) == 1 {
		// Single Redis instance
		client := redis.NewClient(&redis.Options{
			Addr:     hosts[0],
			Password: password, // Password (optional)
		})
		pingTest(client)
		return &RedisHandler{Client: client, IsCluster: false}
	} else {
		// Redis cluster
		cluster := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hosts,
			Password: password, // Password (optional)
		})
		clusterPingTest(cluster)
		return &RedisHandler{Cluster: cluster, IsCluster: true}
	}
}

func (r *RedisHandler) GetKeys(pattern string, db int) []string {
	if r.IsCluster {
		return r.Cluster.Keys(pattern).Val()
	}
	// Select the database
	r.Client.Select(int64(db))
	return r.Client.Keys(pattern).Val()
}

func (r *RedisHandler) MigrateKey(keys []string, destination *RedisHandler) {
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		val, err := r.Client.Dump(key).Result()
		if err != nil {
			fmt.Printf("Failed to dump key %s: %v\n", key, err)
			continue
		}
		err = destination.Client.Restore(key, 0, val).Err()
		if err != nil {
			fmt.Printf("Failed to restore key %s: %v\n", key, err)
		} else {
			fmt.Printf("Migrated key: %s\n", key)
		}
	}
}

func (r *RedisHandler) MigrateAllDatabases(destination *RedisHandler) {
	for db := 0; db < 16; db++ {
		r.Client.Select(int64(db))
		keys := r.Client.Keys("*").Val()
		fmt.Printf("Migrating database %d: %d keys\n", db, len(keys))
		r.MigrateKey(keys, destination)
	}
}

func clusterPingTest(redisClient *redis.ClusterClient) {
	pingTest := redisClient.Ping()
	pingMessage, err := pingTest.Result()
	if err != nil {
		log.Fatalln("Cluster ping failed:", pingMessage)
	}
}

func pingTest(redisClient *redis.Client) {
	pingTest := redisClient.Ping()
	pingMessage, err := pingTest.Result()
	if err != nil {
		log.Fatalln("Ping failed:", pingMessage)
	}
}
