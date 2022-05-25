package jenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"projectXBackend/util"
	"regexp"
	"strings"
)

type Jenkins struct {
	threadDumpPattern *regexp.Regexp
	metricsPattern    *regexp.Regexp
}

func NewJenkinsConnector() *Jenkins {
	//threadDumpPattern, _ := regexp.Compile("(THREAD DUMP FOR TEST FAILURE:((.|\n)*?))(INFO|DEBUG|TRACE|WARN|ERROR)")
	threadDumpPattern, _ := regexp.Compile("(?m)(THREAD DUMP FOR TEST FAILURE.*)(.|\n)*?^\"((.|\n)*?)(^.*INFO|DEBUG|TRACE|WARN|ERROR)")
	metricsPattern, _ := regexp.Compile("Metrics recorded during the test:((.|\\n)*?)Finished Running Test.*")

	return &Jenkins{
		threadDumpPattern: threadDumpPattern,
		metricsPattern:    metricsPattern,
	}
}

func (j *Jenkins) LogParse(logFiles map[string]string, failedTests map[string][]string, filePaths map[string]string, runID string) (map[string]map[string]string, error) {
	// - Log parsing (Zoltan script + thread dump + metric parsing + failure reason)
	logsDetailedMap := make(map[string]map[string]string)

	for className, logContent := range logFiles {
		fmt.Println("PARSE STARTED", className)
		//Metrics
		metricsDump := j.metricsPattern.FindAllString(logContent, 1)

		// Extract test log
		for _, failedTestName := range failedTests[className] {
			logIdentifier := util.CreateIdentifier(className, failedTestName, runID)
			if _, ok := logsDetailedMap[logIdentifier]; !ok {
				logsDetailedMap[logIdentifier] = make(map[string]string)
			}
			cmd := exec.Command("./hztest_zoltan.py", "--file", filePaths[className], "--test", failedTestName)
			out, err := cmd.Output()
			os.Remove(filePaths[className])

			if err != nil {
				println("ZOLTAN", err)
				return nil, err
			}

			logsDetailedMap[logIdentifier]["log"] = string(out)
			if len(metricsDump) != 0 {
				logsDetailedMap[logIdentifier]["metrics"] = metricsDump[0]
			}
			AppendOrCreateIssue(className+"#"+failedTestName, "http://localhost:3000")

		}
		// threadDump
		threadDumps := j.extractThreadDump(logContent, failedTests[className])
		for testName, threadDump := range threadDumps {
			identifier := util.CreateIdentifier(className, testName, runID)
			logsDetailedMap[identifier]["thread_dump"] = threadDump
		}

	}
	return logsDetailedMap, nil
}

func (j *Jenkins) extractThreadDump(logFile string, failedTests []string) map[string]string {
	threadDumps := make([]string, 20)
	threadDumpsMatches := j.threadDumpPattern.FindAllStringSubmatch(logFile, 20)
	for _, threadDumpParts := range threadDumpsMatches {
		threadDump := ""
		for _, threadDumpPart := range threadDumpParts[1 : len(threadDumpParts)-1] {
			threadDump += threadDumpPart
		}
		strings.TrimSpace(threadDump)
		threadDumps = append(threadDumps, threadDump)
	}
	threadDumpMap := make(map[string]string)
	for _, threadDump := range threadDumps {
		for _, testName := range failedTests {
			if strings.Contains(threadDump, testName) {
				idx := strings.LastIndex(threadDump, "\n")
				threadDumpMap[testName] = threadDump[:idx]
				continue
			}
		}
	}
	return threadDumpMap
}

type Payload struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Labels    []string `json:"labels"`
}

type Issue struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Comment struct {
	Body string `json:"body"`
}

func AppendOrCreateIssue(testName string, link string) {
	url := "https://api.github.com/repos/ramizdundar/hazelcast/issues"
	token := os.Getenv("BESTEST_TOKEN")
	fmt.Println("URL:>", url)

	req1, _ := http.NewRequest("GET", url, nil)
	req1.Header.Set("Accept", "application/vnd.github.v3+json")
	req1.SetBasicAuth("ramizdundar", token)
	resp1, _ := http.DefaultClient.Do(req1)

	body1, _ := ioutil.ReadAll(resp1.Body)

	var issues []Issue
	json.Unmarshal(body1, &issues)

	defer resp1.Body.Close()

	for _, issue := range issues {
		if issue.Title == testName {
			comment := Comment{
				Body: link,
			}
			commentBytes, _ := json.Marshal(comment)

			commentBody := bytes.NewReader(commentBytes)
			req3, _ := http.NewRequest("POST", issue.Url+"/comments", commentBody)

			req3.Header.Set("Accept", "application/vnd.github.v3+json")
			req3.SetBasicAuth("ramizdundar", token)

			resp3, _ := http.DefaultClient.Do(req3)

			defer resp3.Body.Close()
			return
		}
	}

	data := Payload{
		Title:     testName,
		Body:      link,
		Assignees: []string{"ramizdundar"},
		Labels:    []string{"bug"},
	}

	payloadBytes, _ := json.Marshal(data)

	body := bytes.NewReader(payloadBytes)

	req2, _ := http.NewRequest("POST", url, body)
	req2.Header.Set("Accept", "application/vnd.github.v3+json")
	req2.SetBasicAuth("ramizdundar", token)

	resp2, _ := http.DefaultClient.Do(req2)
	defer resp2.Body.Close()
}
