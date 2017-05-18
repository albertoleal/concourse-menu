package menu

import (
	"fmt"

	"github.com/albertoleal/concourse-menu/fly"
	"github.com/concourse/atc"
	"github.com/concourse/go-concourse/concourse"
)

type Pipeline struct {
	client fly.Client
	name   string
}

type Build struct {
	atc.Build

	ConcourseURL string
}

func NewPipeline(client fly.Client, name string) *Pipeline {
	return &Pipeline{
		client: client,
		name:   name,
	}
}

func (p *Pipeline) BuildsForGroup(name string) ([]Build, error) {
	builds := []Build{}

	pipeline, err := p.client.Pipeline(p.name)
	if err != nil {
		return []Build{}, err
	}

	team := p.client.Team()
	for _, group := range pipeline.Groups {

		if group.Name == name {
			for _, job := range group.Jobs {
				lastBuilds, _, _, err := team.JobBuilds(p.name, job, concourse.Page{Limit: 1})
				if err != nil {
					return []Build{}, err
				}

				if len(lastBuilds) > 0 {
					concourseURL := fmt.Sprintf("%s%s", p.client.ConcourseURL(), lastBuilds[0].URL)
					b := Build{lastBuilds[0], concourseURL}
					builds = append(builds, b)
				}
			}
		}

	}

	return builds, nil
}
