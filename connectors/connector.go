package connectors

import (
	j "projectXBackend/connectors/jenkins"
)

type Connector interface {
	LogParse(logFiles map[string]string, failedTests map[string][]string, runID string) (map[string]map[string]string, error)
}

func NewConnector(connectorType string) Connector {
	switch connectorType {
	case "jenkins":
		return j.NewJenkinsConnector()
		//case "github-actions":
		//	return g.NewGithubActionsConnector()
	}
	return nil
}
