package jenkins

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Jenkins struct {
	threadDumpPattern *regexp.Regexp
	metricsPattern    *regexp.Regexp
}

func NewJenkinsConnector() *Jenkins {
	dumpPattern, _ := regexp.Compile("(THREAD DUMP FOR TEST FAILURE:((.|\n)*?))(INFO|DEBUG|TRACE|WARN|ERROR)")
	metricsPattern, _ := regexp.Compile("Metrics recorded during the test:((.|\\n)*?)Finished Running Test.*")

	return &Jenkins{
		threadDumpPattern: dumpPattern,
		metricsPattern:    metricsPattern,
	}
}

func (j *Jenkins) LogParse(logFiles map[string]string, failedTests map[string][]string, filePaths map[string]string, runID string) (map[string]map[string]string, error) {
	// - Log parsing (Zoltan script + thread dump + metric parsing + failure reason)
	logsDetailedMap := make(map[string]map[string]string)

	for fileName, logContent := range logFiles {
		fmt.Println("PARSE STARTED", fileName)
		//Metrics
		metricsDump := j.metricsPattern.FindAllString(logContent, 1)

		// Extract test log
		for _, failedTestName := range failedTests[fileName] {
			logIdentifier := failedTestName + "_" + runID
			if _, ok := logsDetailedMap[logIdentifier]; !ok {
				logsDetailedMap[logIdentifier] = make(map[string]string)
			}
			cmd := exec.Command("./hztest_zoltan.py", "--file", filePaths[fileName], "--test", failedTestName)
			out, err := cmd.Output()
			os.Remove(filePaths[fileName])

			if err != nil {
				println("ZOLTAN", err)
				return nil, err
			}

			logsDetailedMap[logIdentifier]["log"] = string(out)
			if len(metricsDump) != 0 {
				logsDetailedMap[logIdentifier]["metrics"] = metricsDump[0]
			}

		}
		// threadDump
		threadDumps := j.extractThreadDump(logContent, failedTests[fileName])
		for testName, threadDump := range threadDumps {
			logIdentifier := testName + "_" + runID
			logsDetailedMap[logIdentifier]["thread_dump"] = threadDump
		}

	}
	return logsDetailedMap, nil
}

func (j *Jenkins) extractThreadDump(logFile string, failedTests []string) map[string]string {
	threadDumps := j.threadDumpPattern.FindAllString(logFile, 20)
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
