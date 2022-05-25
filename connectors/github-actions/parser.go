package github_actions

type GithubActions struct {
}

func NewGithubActionsConnector() *GithubActions {
	return &GithubActions{}
}

func (j *GithubActions) LogParse(logFiles map[string]string, failedTests map[string][]string, runID string) (map[string]map[string]string, error) {
	return nil, nil
}
