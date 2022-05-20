package hazelcast

import (
	"context"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/types"
	"time"
)

const (
	LogFileMap  = "log-file-map"
	MetaDataMap = "log-map"
)

func GetHazelcastClient() *hazelcast.Client {
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
		fmt.Println(err)
		panic(err)
	}
	_, err = client.GetMap(ctx, LogFileMap)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = client.GetMap(ctx, MetaDataMap)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return client
}
