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

// PullsAdd creates a new pull request on the specified repository.
func PullsAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner               string `json:"owner"  zeit:"required"`
			Repo                string `json:"repo"   zeit:"required"`
			Title               string `json:"title"  zeit:"required"`
			Head                string `json:"head"   zeit:"required"`
			Base                string `json:"base"   zeit:"required"`
			Body                string `json:"body"`
			Issue               int    `json:"issue"`
			MaintainerCanModify bool   `json:"maintainer_can_modify"`
			Draft               bool   `json:"draft"`
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
	if t == "" {
		c.WriteJSON(http.StatusUnauthorized, cloudpkgs.H{"status": http.StatusUnauthorized})
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	opt := &gogithub.NewPullRequest{
		Title:               &github.Title,
		Head:                &github.Head,
		Base:                &github.Base,
		MaintainerCanModify: &github.MaintainerCanModify,
		Draft:               &github.Draft,
	}

	pull, _, err := client.Client.PullRequests.Create(ctx, github.Owner, github.Repo, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request": pull})
}
