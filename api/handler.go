package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
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
	//logMap, _ := a.Hz.GetMap(ctx, hz.LogFileMap)
	////metadataMap, _ := a.Hz.GetMap(ctx, hz.MetaDataMap)
	//values, err := logMap.GetValues(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//for _, a := range values {
	//	fmt.Println(a.(map[string]string))
	//}

	//testMap, _ := a.Hz.GetMap(ctx, hz.TestMap)
	//test1RunIDs := []string{"a", "b", "c", "d"}
	//test2RunIDs := []string{"e", "f"}
	//testMap.Put(ctx, "test1", test1RunIDs)
	//testMap.Put(ctx, "test2", test2RunIDs)
	//tests, _ := testMap.GetKeySet(ctx)
	//fmt.Println(tests)
	logMap, _ := a.Hz.GetMap(ctx, hz.LogMap)
	//logs := make(map[string]string)
	//logs["extracted_test_log"] = "Logslogslogslogslogs"
	//fmt.Println(logs)
	//logMap.Put(ctx, "test1_a", logs)
	size, _ := logMap.Size(ctx)
	fmt.Println(size)
	var kk map[string]string
	aaa, err := logMap.Get(ctx, "test1_a")
	fmt.Println(err)
	kk = aaa.(map[string]string)
	//fmt.Println(err)
	fmt.Println(kk)

}

func (a *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := struct {
		Metadata    map[string]string   `json:"metadata"`
		FailedTests map[string][]string `json:"failed_tests"` // filename -> file content
		Logfiles    map[string]string   `json:"log_files,omitempty"`
	}{}

	if r.Header.Get("Content-type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Incorrect JSON body in request: %v\n", err)
			return
		}
	} else {
		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid multipart form request: %v\n", err)

			return
		}
		err = json.Unmarshal([]byte(r.Form["req"][0]), &req)
		if err != nil {
			fmt.Fprintf(w, "Incorrect JSON body in the form: %v\n", err)
			return
		}
		logFiles := make(map[string]string)
		for filename, testNames := range req.FailedTests {
			// Get file from Form
			file, _, err := r.FormFile(filename)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to get log file: %v\n", err)
				return
			}
			// Read entire file content, giving us little control but
			// making it very simple. No need to close the file.
			filebytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Convert []byte to string and print to screen
			fileContent := string(filebytes)
			logFiles[filename] = fileContent

			a.Hz.AppendTestRunID(ctx, testNames, req.Metadata["runID"])
		}
		req.Logfiles = logFiles
	}

	con := req.Metadata["connector"]
	connector := cn.NewConnector(con)

	parsedTestLogs, err := connector.LogParse(req.Logfiles, req.FailedTests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to parse log files: %v\n", err)
		return
	}

	for k, v := range parsedTestLogs {
		for k1, v1 := range v {
			fmt.Println(k, k1, v1[0:100])
		}
	}

	w.WriteHeader(200)
	return
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tests, err := a.Hz.GetTestList(ctx)
	if err != nil {
		log.Print("HZ GetTestList: ", err) // log it
		return
	}
	a.createPage(w, "homepage", tests)
}

func (a *App) testRunIDs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	testName, _ := mux.Vars(r)["test-name"]
	testRunIDs, err := a.Hz.GetTestRunIDs(ctx, testName)
	if err != nil {
		log.Print("HZ GetTestRunIDs: ", err) // log it
		return
	}
	pageData := struct {
		TestName string
		RunIDs   interface{}
	}{
		TestName: testName,
		RunIDs:   testRunIDs,
	}
	a.createPage(w, "testrun", pageData)
}

func (a *App) testLogs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logIdentifier, _ := mux.Vars(r)["log-identifier"]

	logs, err := a.Hz.GetLogs(ctx, logIdentifier)
	if err != nil {
		log.Print("HZ GetLogs: ", err) // log it
		return
	}

	logMap := logs.(map[string]string)
	logNames := make([]string, 0, len(logMap))
	for k := range logMap {
		logNames = append(logNames, k)
	}

	pageData := struct {
		LogIdentifier string
		LogNames      interface{}
	}{
		LogIdentifier: logIdentifier,
		LogNames:      logNames,
	}
	a.createPage(w, "logs", pageData)
}

func (a *App) testLogDetail(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logIdentifier, _ := mux.Vars(r)["log-identifier"]
	logName, _ := mux.Vars(r)["log-name"]

	logs, err := a.Hz.GetLogs(ctx, logIdentifier)
	if err != nil {
		log.Print("HZ GetLogs: ", err) // log it
		return
	}

	logMap := logs.(map[string]string)
	logDetail := logMap[logName]

	pageData := struct {
		LogIdentifier string
		LogName       string
		LogDetail     string
	}{
		LogIdentifier: logIdentifier,
		LogName:       logName,
		LogDetail:     logDetail,
	}
	a.createPage(w, "logDetail", pageData)
}

func (a *App) createPage(w http.ResponseWriter, page string, data interface{}) {
	t, err := template.ParseFiles("template/" + page + ".html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, data) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
