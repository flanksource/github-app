package runners

import (
	"context"
	"fmt"
	"github.com/flanksource/github-app/config"
	"github.com/flanksource/github-app/pkg/util"
	"github.com/google/go-github/v32/github"
)

type  GhUtil struct {
	cfg *config.Config
}

func (h *GhUtil)  Cleanup() (bool, *[]error) {
	haveDeleted := false
	var errors *[]error

	client := util.GetPatClient(h.cfg.Secrets.GhPat)

	runners, _, err := client.Actions.ListRunners(context.TODO(), h.cfg.Runners.Owner, h.cfg.Runners.Repo,&github.ListOptions{})
	if err != nil {
		errorStorage := make([]error, 1)
		errors = &errorStorage
		*errors = append(*errors, err)
		return false, errors
	}
	for _, runner := range runners.Runners {
		if *runner.Status=="offline" {
			rsp, err := client.Actions.RemoveRunner(context.TODO(), h.cfg.Runners.Owner+"nope", h.cfg.Runners.Repo,*runner.ID)
			fmt.Printf("%v",rsp)
			if err != nil{
				if errors == nil {
					errorStorage := make([]error, 1)
					errors = &errorStorage
					*errors = make([]error, 1)
				}
				*errors = append(*errors, err)
				continue
			}
		}
	}

	return haveDeleted, errors

}
