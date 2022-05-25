package hazelcast

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/types"
	"time"
)

const (
	LogMap      = "log"
	MetadataMap = "metadata"
	TestMap     = "tests"
)

type HZ struct {
	*hazelcast.Client
}

func GetHazelcastClient() *HZ {
	gob.Register(map[string]string{})
	gob.Register(map[string]bool{})
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
	return &HZ{client}
}

func (hz *HZ) GetTestList(ctx context.Context) ([]interface{}, error) {
	testMap, _ := hz.GetMap(ctx, TestMap)

	tests, err := testMap.GetKeySet(ctx)
	if err != nil {
		//return nil, fmt.Errorf("failed to get keys from map %s: %v", err)
		return nil, err
	}
	return tests, nil
}

func (hz *HZ) GetTestRunIDs(ctx context.Context, testName string) interface{} {
	testMap, _ := hz.GetMap(ctx, TestMap)

	testRunIDs, _ := testMap.Get(ctx, testName)
	return testRunIDs

}

func (hz *HZ) GetLogs(ctx context.Context, logIdentifier string) (interface{}, error) {
	logMap, _ := hz.GetMap(ctx, LogMap)

	logs, err := logMap.Get(ctx, logIdentifier)
	if err != nil {
		//return nil, fmt.Errorf("failed to get keys from map %s: %v", err)
		return nil, nil
	}
	return logs, nil
}

func (hz *HZ) AppendTestRunID(ctx context.Context, testNames []string, testRunID string) {
	testMap, _ := hz.GetMap(ctx, TestMap)

	for _, testName := range testNames {
		testRunIDs, _ := testMap.Get(ctx, testName)
		if testRunIDs == nil {
			_, err := testMap.Put(ctx, testName, []string{testRunID})
			if err != nil {
				fmt.Println("AppendTestRunID Create", err)
				return
			}
			continue
		}

		_, err := testMap.Put(ctx, testName, append(testRunIDs.([]string), testRunID))
		if err != nil {
			fmt.Println("AppendTestRunID", err)
			return
		}
	}
}

//func (hz *HZ) StoreMetadata(ctx context.Context, metadata testNames []string, ) (interface{}, error) {
//	logMap, _ := hz.GetMap(ctx, LogMap)
//
//	logs, err := logMap.Get(ctx, logIdentifier)
//	if err != nil {
//		//return nil, fmt.Errorf("failed to get keys from map %s: %v", err)
//		return nil, nil
//	}
//	return logs, nil
//}

//Metadata struct {
//RunID          string      `json:"runID"`
//NodeId         interface{} `json:"nodeId"`
//CommitId       string      `json:"commitId"`
//JenkinsJobName string      `json:"jenkinsJobName"`
//GitRepoUrl     string      `json:"gitRepoUrl"`
//Connector      string      `json:"connector"`
//} `json:"metadata"`
