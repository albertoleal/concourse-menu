package menu

import (
	"fmt"
	"time"
)

const (
	BUILDING = "Building"
	SUCCESS  = "Success"
	FAILURE  = "Failure"
)

func ConvertBuildsToProject(concourseHost string, pipelineGroupName string, builds []Build) string {
	var (
		status        string
		lastBuildTime time.Time
		failed        string
		running       string
		buildURL      string
		lastBuildURL  string
	)

	for _, build := range builds {
		endTime := time.Unix(build.EndTime, 0)
		if build.EndTime > lastBuildTime.Unix() {
			lastBuildTime = endTime.Local()
			lastBuildURL = build.URL
		}

		if build.IsRunning() {
			running = build.URL
		}

		switch build.Status {
		case "aborted", "errored", "failed":
			failed = build.URL
		case "succeeded":
			status = SUCCESS
		}
	}

	if running != "" {
		status = BUILDING
		buildURL = running
	} else if failed != "" {
		status = FAILURE
		buildURL = failed
	} else {
		buildURL = lastBuildURL
	}

	concourseURL := fmt.Sprintf("%s%s", concourseHost, buildURL)
	return fmt.Sprintf("<Project name=\"%s\" activity=\"%s\" lastBuildStatus=\"%s\" lastBuildTime=\"%s\" webUrl=\"%s\" \n/>", pipelineGroupName, status, status, lastBuildTime.Format("2006-01-02T15:04:05Z07:00"), concourseURL)
}
