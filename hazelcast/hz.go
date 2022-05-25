package hazelcast

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/types"
	"projectXBackend/util"
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
		cTestName, _ := testMap.Get(ctx, className)
		if cTestName == nil {
			testMap.Put(ctx, className, testNames)
			continue
		}
		testMap.Put(ctx, className, removeDuplicateStr(append(testNames, cTestName.([]string)...)))
	}
}

func (hz *HZ) GetTestNames(ctx context.Context) ([]types.Entry, error) {
	testMap, _ := hz.GetMap(ctx, TestMap)
	classNames, _ := testMap.GetKeySet(ctx)
	tests, _ := testMap.GetAll(ctx, classNames...)

	return tests, nil
}

func (hz *HZ) StoreTestRunID(ctx context.Context, testNames map[string][]string, runID string) {
	testRunIdMap, _ := hz.GetMap(ctx, RunIDMap)

	for className, testNames := range testNames {
		for _, testName := range testNames {
			runIdIdentifier := util.CreateIdentifier(className, testName)
			fmt.Println("AAA", runIdIdentifier)
			runIds, _ := testRunIdMap.Get(ctx, runIdIdentifier)
			if runIds == nil {
				testRunIdMap.Put(ctx, runIdIdentifier, []string{runID})
				continue
			}
			testRunIdMap.Put(ctx, runIdIdentifier, removeDuplicateStr(append([]string{runID}, runIds.([]string)...)))
		}
	}
}

func (hz *HZ) GetTestRunIDs(ctx context.Context, className string, testName string) interface{} {
	testRunIdMap, _ := hz.GetMap(ctx, RunIDMap)
	size, err := testRunIdMap.Size(ctx)
	if err != nil {
		return nil
	}
	fmt.Println(size)
	aa, err := testRunIdMap.GetKeySet(ctx)
	if err != nil {
		return nil
	}
	fmt.Println(aa)
	runIdIdentifier := util.CreateIdentifier(className, testName)
	fmt.Println("AAA", runIdIdentifier)

	testRunIDs, _ := testRunIdMap.Get(ctx, runIdIdentifier)
	return testRunIDs
}

func (hz *HZ) StoreLogs(ctx context.Context, parsedTestLogs map[string]map[string]string) {
	logMap, _ := hz.GetMap(ctx, LogMap)

	for identifier, logs := range parsedTestLogs {
		_, _ = logMap.Put(ctx, identifier, logs)
	}
}

func (hz *HZ) GetLogs(ctx context.Context, className string, testName string, runID string) (interface{}, error) {
	logMap, _ := hz.GetMap(ctx, LogMap)

	logs, err := logMap.Get(ctx, util.CreateIdentifier(className, testName, runID))
	if err != nil {
		//return nil, fmt.Errorf("failed to get keys from map %s: %v", err)
		return nil, nil
	}
	return logs, nil
}

//func (hz *HZ) StoreMetadata(ctx context.Context, metadata map[string]string, testNames []string) {
//	metadataMap, _ := hz.GetMap(ctx, MetadataMap)
//
//	for _, testName := range testNames {
//		logIdentifier := testName + "_" + metadata["runID"]
//		metadataMap.Put(ctx, logIdentifier, metadata)
//	}
//}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
