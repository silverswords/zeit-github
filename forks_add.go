package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// ForksAdd  creates a fork of the specified repository.
// This method might return an *AcceptedError and a status code of
// 202. This is because this is the status that GitHub returns to signify that
// it is now computing creating the fork in a background task. In this event,
// the Repository value will be returned, which includes the details about the pending fork.
// A follow up request, after a delay of a second or so, should result
// in a successful request.
func ForksAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner        string `json:"owner" zeit:"required"`
			Repo         string `json:"repo"  zeit:"required"`
			Organization string `json:"organization"`
		}
	)

	c := cloudpkgs.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusBadRequest, cloudpkgs.H{"status": http.StatusBadRequest})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusPreconditionRequired, cloudpkgs.H{"status": http.StatusPreconditionRequired})
		return
	}

	token := c.Request.Header
	t := token.Get("Authorization")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	opt := &gogithub.RepositoryCreateForkOptions{
		Organization: github.Organization,
	}

	gist, _, err := client.Client.Repositories.CreateFork(ctx, github.Owner, github.Repo, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "gists": gist})
}
