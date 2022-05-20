package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	cn "projectXBackend/connectors"
	hz "projectXBackend/hazelcast"
)

func (a *App) testHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	//m, err := a.Hz.GetMap(ctx, hz.MetaDataMap)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//m.Clear(ctx)
	//m1 := make(map[string]string)
	//m1["try"] = "koydukoldu"
	//m1["allahallah"] = "sasirtti"
	//m.Put(ctx, "aa", "aaa")
	//m.Put(ctx, 3, m1)
	//
	//v, _ := m.Get(ctx, "aa")
	//v2, _ := m.Get(ctx, 3)
	//fmt.Println(v)
	//fmt.Println(v2)
	//size, _ := m.Size(ctx)
	//fmt.Println(size)
	//w.Write([]byte(v2.(map[string]string)["try"]))
	logMap, _ := a.Hz.GetMap(ctx, hz.LogFileMap)
	//metadataMap, _ := a.Hz.GetMap(ctx, hz.MetaDataMap)
	values, err := logMap.GetValues(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, a := range values {
		fmt.Println(a.(map[string]string))
	}
}

//type Metadata struct {
//	CommitID  string `json:"commit_id"`
//	Timestamp string `json:"timestamp"`
//}

func (a *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	req := struct { // TODO: add validator
		Metadata    map[string]string   `json:"metadata"`
		Logfiles    map[string]string   `json:"log_files"`    // filename -> file content
		FailedTests map[string][]string `json:"failed_tests"` // filename -> file content
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Incorrect JSON body in request: %v\n", err)
		return
	}

	//con := req.Metadata["connector"]
	con := "jenkins"
	connector := cn.NewConnector(con)

	parsedTestLogs, err := connector.LogParse(req.Logfiles, req.FailedTests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to parse log files: %v\n", err)
		return
	}
	ctx := context.Background()

	logMap, _ := a.Hz.GetMap(ctx, hz.LogFileMap)
	logMap.Put(ctx, req.Metadata["id"], parsedTestLogs)

	metadataMap, _ := a.Hz.GetMap(ctx, hz.MetaDataMap)
	metadataMap.Put(ctx, req.Metadata["id"], req.Metadata)
}

func (a *App) getLogsHandler(w http.ResponseWriter, r *http.Request) {

}
