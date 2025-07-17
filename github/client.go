package github

import (
	"encoding/json"
	"fmt"

	"github.com/cli/go-gh"
)

func getStatusIcon(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return "✅"
		case "failure":
			return "❌"
		case "cancelled":
			return "🚫"
		case "skipped":
			return "⏭️"
		default:
			return conclusion
		}
	case "in_progress":
		return "🔄"
	case "queued":
		return "⏳"
	default:
		return status
	}
}

type Workflow struct {
	Name  string `json:"name"`
	State string `json:"state"`
	ID    int    `json:"id"`
}

func (w Workflow) FilterValue() string { return w.Name }
func (w Workflow) Title() string       { return w.Name }
func (w Workflow) Description() string { return w.State }

type WorkflowRun struct {
	ID          int    `json:"databaseId"`
	Number      int    `json:"number"`
	DisplayTitle string `json:"displayTitle"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	Event       string `json:"event"`
	HeadBranch  string `json:"headBranch"`
	URL         string `json:"url"`
}

func (r WorkflowRun) FilterValue() string { return r.DisplayTitle }
func (r WorkflowRun) Title() string {
	return fmt.Sprintf("%s #%d %s", getStatusIcon(r.Status, r.Conclusion), r.Number, r.DisplayTitle)
}
func (r WorkflowRun) Description() string {
	return fmt.Sprintf("%s • %s • %s", r.HeadBranch, r.Event, r.URL)
}

type Job struct {
	ID         int    `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	URL        string `json:"url"`
	Steps      []Step `json:"steps"`
}

func (j Job) FilterValue() string { return j.Name }
func (j Job) Title() string       { return fmt.Sprintf("%s %s", getStatusIcon(j.Status, j.Conclusion), j.Name) }
func (j Job) Description() string { return j.URL }

type Step struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Number     int    `json:"number"`
}

// Client provides methods for interacting with the GitHub API.
type Client struct{}

// NewClient creates a new GitHub client.
func NewClient() *Client {
	return &Client{}
}

// ListWorkflows fetches a list of GitHub Actions workflows for the specified repository.
func (c *Client) ListWorkflows(repo string) ([]Workflow, error) {
	args := []string{"workflow", "list", "--json", "name,state,id"}
	if repo != "" {
		args = append(args, "--repo", repo)
	}
	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	var workflows []Workflow
	err = json.Unmarshal(stdout.Bytes(), &workflows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow list JSON: %w", err)
	}

	return workflows, nil
}

// ListWorkflowRuns fetches a list of runs for a specific workflow.
func (c *Client) ListWorkflowRuns(repo string, workflowID int) ([]WorkflowRun, error) {
	args := []string{
		"run", "list",
		"--workflow", fmt.Sprintf("%d", workflowID),
		"--json", "databaseId,number,displayTitle,status,conclusion,event,headBranch,url",
	}
	if repo != "" {
		args = append(args, "--repo", repo)
	}

	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow runs: %w", err)
	}

	var runs []WorkflowRun
	err = json.Unmarshal(stdout.Bytes(), &runs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow runs JSON: %w", err)
	}

	return runs, nil
}

// ListJobs fetches a list of jobs for a specific workflow run.
func (c *Client) ListJobs(repo string, runID int) ([]Job, error) {
	args := []string{
		"run", "view", fmt.Sprintf("%d", runID),
		"--json", "jobs",
	}
	if repo != "" {
		args = append(args, "--repo", repo)
	}

	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	var result struct {
		Jobs []Job `json:"jobs"`
	}
	err = json.Unmarshal(stdout.Bytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jobs JSON: %w", err)
	}

	return result.Jobs, nil
}

// GetJobLogs fetches the logs for a specific job.
func (c *Client) GetJobLogs(repo string, jobID int) (string, error) {
	// Use GitHub CLI to get logs directly (fixed in gh CLI 2.75+)
	args := []string{
		"run", "view", "--log",
		"--job", fmt.Sprintf("%d", jobID),
	}
	if repo != "" {
		args = append(args, "--repo", repo)
	}

	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get job logs: %w", err)
	}

	return stdout.String(), nil
}

