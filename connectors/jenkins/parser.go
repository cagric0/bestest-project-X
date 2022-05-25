package jenkins

import (
	"fmt"
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

func (j *Jenkins) LogParse(logFiles map[string]string, failedTests map[string][]string, runID string) (map[string]map[string]string, error) {
	// - Log parsing (Zoltan script + thread dump + metric parsing + failure reason)
	logsDetailedMap := make(map[string]map[string]string)

	for className, logContent := range logFiles {
		fmt.Println("PARSE STARTED", className)
		//Metrics
		metricsDump := j.metricsPattern.FindAllString(logContent, 1)

		// Extract test log
		for _, failedTestName := range failedTests[className] {
			logIdentifier := failedTestName + "_" + runID
			if _, ok := logsDetailedMap[logIdentifier]; !ok {
				logsDetailedMap[logIdentifier] = make(map[string]string)
			}
			//fmt.Println("./hztest_zoltan.py", "--file", logContent, "--test", failedTestName)
			cmd := exec.Command("./hztest_zoltan.py", "--file", logContent, "--test", failedTestName)
			out, err := cmd.Output()
			if err != nil {
				println(err.Error())
				return nil, err
			}

			logsDetailedMap[logIdentifier]["log"] = string(out)
			logsDetailedMap[logIdentifier]["metrics"] = metricsDump[0]

		}
		// threadDump
		threadDumps := j.extractThreadDump(logContent, failedTests[className])
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
				fmt.Println("THREAD DUMP", testName, threadDump[0:200])
				continue
			}
		}
	}
	return threadDumpMap
}
