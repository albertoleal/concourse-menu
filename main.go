package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/albertoleal/concourse-menu/fly"
	"github.com/albertoleal/concourse-menu/menu"
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
	client   = fly.NewClient(target, username, password, team)
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	rmq_1_7_release := createProject("rabbitmq-1.7", "RELEASE", "RabbitMQ-1.7: Release")
	rmq_1_7_product := createProject("rabbitmq-1.7", "PRODUCT", "RabbitMQ-1.7: Product")
	rmq_1_7_cleanup := createProject("rabbitmq-1.7", "CLEANUP", "RabbitMQ-1.7: Cleanup")

	rmq_1_8_release := createProject("rabbitmq-1.8", "RELEASE", "RabbitMQ-1.8: Release")
	rmq_1_8_product := createProject("rabbitmq-1.8", "PRODUCT", "RabbitMQ-1.8: Product")
	rmq_1_8_cleanup := createProject("rabbitmq-1.8", "CLEANUP", "RabbitMQ-1.8: Cleanup")

	rmq_1_9_release := createProject("rabbitmq", "RELEASE", "RabbitMQ-1.9: Release")
	rmq_1_9_product := createProject("rabbitmq", "PRODUCT", "RabbitMQ-1.9: Product")
	rmq_1_9_cleanup := createProject("rabbitmq", "CLEANUP", "RabbitMQ-1.9: Cleanup")

	cf_rmq_release_1_8 := createProject("rabbitmq-1.8", "CF-RABBITMQ-SERVER", "RabbitMQ-1.8: CF RMQ Server")
	cf_rmq_release_1_9 := createProject("rabbitmq", "CF-RABBITMQ-SERVER", "RabbitMQ-1.9: CF RMQ Server")

	cf_rmq_broker_release_1_9 := createProject("rabbitmq", "MULTITENANT-BROKER", "RabbitMQ-1.9: Multi-tenant broker")
	cf_rmq_metrics_1_9 := createProject("rabbitmq", "RABBITMQ-METRICS", "RabbitMQ-1.9: Metrics")

	out := fmt.Sprintf("<Projects>%s%s%s%s%s%s%s%s%s%s%s%s%s</Projects>", rmq_1_7_release, rmq_1_7_product, rmq_1_7_cleanup, rmq_1_8_release, rmq_1_8_product, rmq_1_8_cleanup, rmq_1_9_release, rmq_1_9_product, rmq_1_9_cleanup, cf_rmq_release_1_9, cf_rmq_release_1_8, cf_rmq_broker_release_1_9, cf_rmq_metrics_1_9)
	io.WriteString(w, out)
}

func createProject(pipelineName string, group string, name string) string {
	pipeline := menu.NewPipeline(client, pipelineName)
	builds, err := pipeline.BuildsForGroup(group)
	if err != nil {
		panic(err)
	}
	return menu.ConvertBuildsToProject(target, name, builds)
}
