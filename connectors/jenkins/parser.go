package jenkins

import (
	"fmt"
	"os/exec"
)

type Jenkins struct {
}

func NewJenkinsConnector() *Jenkins {
	return &Jenkins{}
}

func (j *Jenkins) LogParse(logFiles map[string]string, failedTests map[string][]string) (map[string]map[string]string, error) {
	// - Log parsing (Zoltan script + thread dump + metric parsing + failure reason)
	logsDetailedMap := make(map[string]map[string]string)
	for className, logContent := range logFiles {
		fmt.Println("PARSE STARTED", className)
		// Extract test log
		for _, failedTestName := range failedTests[className] {
			cmd := exec.Command("/Users/cagriciftci/Desktop/bestest-project-X/hztest_zoltan.py", "--file", logContent, "--test", failedTestName)
			out, err := cmd.Output()
			if err != nil {
				println(err.Error())
				return nil, err
			}
			logsDetailedMap[failedTestName]["extracted_test_log"] = string(out)

			// fmt.Println(string(out))

			// threadDump
			// logsDetailedMap[failedTestName]["thread_dump"] = "threadDump_result"
		}

	}
	return logsDetailedMap, nil
}
