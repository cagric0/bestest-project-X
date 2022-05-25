package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hazelcast/hazelcast-go-client/types"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	cn "projectXBackend/connectors"
	hz "projectXBackend/hazelcast"
)

type Metadata struct {
	RunID          string `json:"runID"`
	NodeId         string `json:"nodeId"`
	CommitId       string `json:"commitId"`
	JenkinsJobName string `json:"jenkinsJobName"`
	GitRepoUrl     string `json:"gitRepoUrl"`
	Connector      string `json:"connector"`
}

func (a *App) clearHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	ctx := context.Background()
	m, err := a.Hz.GetMap(ctx, hz.LogMap)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m1, err := a.Hz.GetMap(ctx, hz.MetadataMap)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m2, err := a.Hz.GetMap(ctx, hz.TestMap)
	m.Clear(ctx)
	m1.Clear(ctx)
	m2.Clear(ctx)
}

func (a *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	ctx := context.Background()
	req := struct {
		Metadata    map[string]string   `json:"metadata"`
		FailedTests map[string][]string `json:"failed_tests"` // filename -> file content
		Logfiles    map[string]string   `json:"log_files,omitempty"`
	}{}
	filePaths := make(map[string]string)

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
			file, fileHeader, err := r.FormFile(filename)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to get log file: %v\n", err)
				return
			}
			filePath := fmt.Sprintf("./uploads/%s", fileHeader.Filename)
			if err := a.createFile(file, filePath); err != nil {
				fmt.Println(err)
				return
			}
			// Read entire file content, giving us little control but
			// making it very simple. No need to close the file.
			file, _ = os.Open(filePath)

			filebytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Convert []byte to string and print to screen
			fileContent := string(filebytes)

			logFiles[filename] = fileContent
			filePaths[filename] = filePath
			a.Hz.AppendTestRunID(ctx, testNames, req.Metadata["runID"])
		}
		req.Logfiles = logFiles
	}

	con := req.Metadata["connector"]
	connector := cn.NewConnector(con)
	a.Hz.StoreTestNames(ctx, req.FailedTests)

	//a.Hz.StoreMetadata(ctx, req.Metadata, req.FailedTests)
	runID := req.Metadata["runID"]

	parsedTestLogs, err := connector.LogParse(req.Logfiles, req.FailedTests, filePaths, runID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to parse log files: %v\n", err)
		return
	}

	logMap, _ := a.Hz.GetMap(ctx, hz.LogMap)
	for logIdentifier, logDetailedMap := range parsedTestLogs {
		_, _ = logMap.Put(ctx, logIdentifier, logDetailedMap)
	}

	w.WriteHeader(200)
	return
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	ctx := context.Background()

	tests, _ := a.Hz.GetTestNames(ctx)
	response := struct {
		Failures []types.Entry `json:"failures"`
	}{
		Failures: tests,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func (a *App) testRunIDs(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	ctx := context.Background()
	testName, _ := mux.Vars(r)["test-name"]
	testRunIDs := a.Hz.GetTestRunIDs(ctx, testName)

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
	enableCors(w)
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
	enableCors(w)
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

func (a *App) createFile(file multipart.File, filePath string) error {
	// Create a new file in the uploads directory
	dst, err := os.Create(filePath)
	if err != nil {
		return nil
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
