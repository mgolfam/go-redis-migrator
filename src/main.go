package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/redis.v3"
)

var (
	sourceCluster      *redis.ClusterClient
	destinationCluster *redis.ClusterClient
	sourceHost         *redis.Client
	destinationHost    *redis.Client

	sourceHostsArray      []string
	destinationHostsArray []string

	sourceIsCluster      = false
	destinationIsCluster = false

	sourceClusterConnected      = false
	destinationClusterConnected = false

	keysMigrated int64
)

func main() {
	sourceHosts := flag.String("sourceHosts", "", "Comma-separated list of source Redis servers (host:port).")
	destinationHosts := flag.String("destinationHosts", "", "Comma-separated list of destination Redis servers (host:port).")
	getKeys := flag.Bool("getKeys", false, "Fetch and display keys from the source Redis cluster.")
	copyData := flag.Bool("copyData", false, "Migrate keys to the destination Redis cluster.")
	keyFilePath := flag.String("keyFile", "", "Path to a file containing keys to migrate.")
	keyFilter := flag.String("keyFilter", "*", "Pattern to match keys for migration.")

	flag.Parse()

	if !*getKeys && !*copyData {
		showHelp()
	}

	if *sourceHosts != "" {
		sourceHostsArray = strings.Split(*sourceHosts, ",")
		connectSourceCluster()
	}

	if *destinationHosts != "" {
		destinationHostsArray = strings.Split(*destinationHosts, ",")
		connectDestinationCluster()
	}

	if *getKeys {
		if !sourceClusterConnected {
			log.Fatalln("Please specify a valid source Redis cluster.")
		}
		keys := getSourceKeys(*keyFilter)
		if len(keys) > 0 {
			for _, key := range keys {
				fmt.Println(key)
			}
		} else {
			fmt.Println("No keys found.")
		}
	}

	if *copyData {
		if !sourceClusterConnected || !destinationClusterConnected {
			log.Fatalln("Both source and destination Redis clusters must be specified.")
		}
		if *keyFilePath != "" {
			if *keyFilter != "*" {
				log.Fatalln("Cannot use keyFilter and keyFile together.")
			}
			file, err := os.Open(*keyFilePath)
			if err != nil {
				log.Fatalln("Unable to open key file:", err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				migrateKey(scanner.Text())
			}
		} else {
			keys := getSourceKeys(*keyFilter)
			for _, key := range keys {
				migrateKey(key)
			}
		}
		fmt.Printf("Migrated %d keys.\n", keysMigrated)
	}
}

func connectSourceCluster() {
	if len(sourceHostsArray) == 1 {
		sourceHost = redis.NewClient(&redis.Options{
			Addr: sourceHostsArray[0],
		})
		sourceIsCluster = false
		hostPingTest(sourceHost)
	} else {
		sourceCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: sourceHostsArray,
		})
		sourceIsCluster = true
		clusterPingTest(sourceCluster)
	}
	sourceClusterConnected = true
}

func connectDestinationCluster() {
	if len(destinationHostsArray) == 1 {
		destinationHost = redis.NewClient(&redis.Options{
			Addr: destinationHostsArray[0],
		})
		destinationIsCluster = false
		hostPingTest(destinationHost)
	} else {
		destinationCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: destinationHostsArray,
		})
		destinationIsCluster = true
		clusterPingTest(destinationCluster)
	}
	destinationClusterConnected = true
}

func getSourceKeys(keyFilter string) []string {
	if sourceIsCluster {
		return sourceCluster.Keys(keyFilter).Val()
	}
	return sourceHost.Keys(keyFilter).Val()
}

func migrateKey(key string) {
	keysMigrated++
	var data string
	var ttl time.Duration

	if sourceIsCluster {
		data = sourceCluster.Dump(key).Val()
		ttl = sourceCluster.PTTL(key).Val()
	} else {
		data = sourceHost.Dump(key).Val()
		ttl = sourceHost.PTTL(key).Val()
	}

	if ttl == -1 {
		ttl = 0
	}

	if destinationIsCluster {
		destinationCluster.Restore(key, ttl, data)
	} else {
		destinationHost.Restore(key, ttl, data)
	}
}

func clusterPingTest(redisClient *redis.ClusterClient) {
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalln("Error pinging Redis cluster:", err)
	}
}

func hostPingTest(redisClient *redis.Client) {
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalln("Error pinging Redis host:", err)
	}
}

func showHelp() {
	fmt.Println(`
Usage:
  -sourceHosts: Comma-separated source Redis servers (e.g., 127.0.0.1:6379,127.0.0.1:6380)
  -destinationHosts: Comma-separated destination Redis servers
  -getKeys: Fetch and print keys from the source
  -copyData: Copy keys to the destination Redis
  -keyFile: File containing keys to copy
  -keyFilter: Pattern to match keys to copy (default: *)
	`)
	os.Exit(0)
}
