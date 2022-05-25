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
	RunIDMap    = "runid"
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

func (hz *HZ) StoreTestNames(ctx context.Context, testNames map[string][]string) {
	testMap, _ := hz.GetMap(ctx, TestMap)

	for className, testNames := range testNames {
		testMap.Put(ctx, className, testNames)
	}
}

func (hz *HZ) GetTestNames(ctx context.Context) ([]types.Entry, error) {
	testMap, _ := hz.GetMap(ctx, TestMap)
	classNames, _ := testMap.GetKeySet(ctx)
	tests, _ := testMap.GetAll(ctx, classNames...)

	return tests, nil
}

func (hz *HZ) GetTestRunIDs(ctx context.Context, testName string) interface{} {
	testMap, _ := hz.GetMap(ctx, RunIDMap)

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
	testMap, _ := hz.GetMap(ctx, RunIDMap)

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

func (hz *HZ) StoreMetadata(ctx context.Context, metadata map[string]string, testNames []string) {
	metadataMap, _ := hz.GetMap(ctx, MetadataMap)

	for _, testName := range testNames {
		logIdentifier := testName + "_" + metadata["runID"]
		metadataMap.Put(ctx, logIdentifier, metadata)
	}
}

type T struct {
	Metadata struct {
		RunID          string      `json:"runID"`
		NodeId         interface{} `json:"nodeId"`
		CommitId       string      `json:"commitId"`
		JenkinsJobName string      `json:"jenkinsJobName"`
		GitRepoUrl     string      `json:"gitRepoUrl"`
		Connector      string      `json:"connector"`
	} `json:"metadata"`
	FailedTests struct {
		ComHazelcastExecutorExecutorServiceTestOutputTxt []string `json:"com.hazelcast.executor.ExecutorServiceTest-output.txt"`
		ComHazelcastExecutorSmallClusterTestOutputTxt    []string `json:"com.hazelcast.executor.SmallClusterTest-output.txt"`
	} `json:"failed_tests"`
}
