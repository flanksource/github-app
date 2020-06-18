package handler

import (
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"net/http"
)

type GHRunners struct {
	githubapp.ClientCreator
	Installations githubapp.InstallationsService
	BaseConfig    *baseapp.HTTPConfig
	Sessions      *scs.Manager
}

func (h *GHRunners) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	//installation, err := h.Installations.GetByOwner(ctx, "philipstaffordwood")
	//if err != nil {
	//	return err
	//}

	//c, err := h.ClientCreator.NewAppClient()
	c, err := h.ClientCreator.NewInstallationClient(9772354)

	if err != nil {
		return errors.Wrap(err, "failed to create github client")
	}
	token, _, err := c.Actions.CreateRegistrationToken(ctx, "philipstaffordwood", "karina")
	if err != nil {
		return err
	}
	token.GetToken()
	fmt.Fprintf(w, "{\"echo\": \"%s\"}", token.GetToken())

	return nil
}
