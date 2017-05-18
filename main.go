package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/albertoleal/concourse-metrics/fly"
	"github.com/concourse/go-concourse/concourse"
)

type Project struct {
	Name            string `xml:"name"`
	Activity        string `xml:"activity"`
	LastBuildStatus string `xml:"lastBuildStatus"`
	LastBuildLabel  string `xml:"lastBuildLabel"`
	LastBuildTime   string `xml:"lastBuildTime"`
	WebUrl          string `xml:"webUrl"`
}

var (
	target   = os.Getenv("FLY_TARGET")
	team     = os.Getenv("TARGET_TEAM")
	username = os.Getenv("TARGET_USERNAME")
	password = os.Getenv("TARGET_PASSWORD")
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	client := fly.NewClient(target, username, password, team)
	builds, _, err := client.Builds(concourse.Page{Limit: 20})
	if err != nil {
		panic(err)
	}

	var projects string
	for _, build := range builds {
		concourseURL := fmt.Sprintf("%s%s", client.ConcourseURL(), build.URL)
		var status string
		switch build.Status {
		case "aborted", "errored":
			status = "Exception"
		case "failed":
			status = "Failure"
		case "succeeded":
			status = "Success"
		}

		projects = projects + fmt.Sprintf("<Project name=\"%s\" activity=\"%s\" lastBuildStatus=\"%s\" lastBuildLabel=\"%s\" lastBuildTime=\"2017-05-08T13:18:38.000+0000\" webUrl=\"%s\" \n/>", fmt.Sprintf("%s/%s", build.PipelineName, build.JobName), status, status, status, concourseURL)
	}
	ps := fmt.Sprintf("<Projects>%s</Projects>", projects)
	if err != nil {
		panic(err)
	}

	io.WriteString(w, ps)
}
