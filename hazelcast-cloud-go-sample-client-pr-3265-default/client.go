package main

import (
	"context"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/types"
	"log"
	"math/rand"
	"time"
)

/**
 *
 * This is boilerplate application that configures client to connect Hazelcast Cloud cluster.
 * After successful connection, it puts random entries into the map.
 *
 * See: https://docs.hazelcast.cloud/docs/go-client
 *
 */
func main() {
	ctx := context.Background()
	config := hazelcast.NewConfig()
	config.Cluster.Name = "pr-3265"
	config.Cluster.Network.SSL.Enabled = false
	config.Cluster.Cloud.Enabled = true
	config.Cluster.Cloud.Token = "R4Fqe9GiqhWlAyUENWJqfHk37z9c9Uaroh3MVxMOA0jpdxCTqa"
	config.Stats.Enabled = true
	config.Stats.Period = types.Duration(time.Second)

	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	defer client.Shutdown(ctx)
	log.Println("Connection Successful!")
	log.Println("Now, `map` will be filled with random entries.")

	mp, err := client.GetMap(ctx, "map")

	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	iterationCounter := 0
	for {
		randKey := string(rune(rand.Intn(100000)))
		_, err := mp.Put(ctx, "key"+randKey, "value"+randKey)
		if err != nil {
			panic(err)
		}
		if iterationCounter++; iterationCounter == 10 {
			iterationCounter = 0
			size, err := mp.Size(ctx)
			if err != nil {
				panic(err)
			}
			log.Println(fmt.Sprintf("Map size: %d", size))
		}
	}
}
