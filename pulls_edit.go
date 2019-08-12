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

// PullsEdit edits a pull request.  pull must not be nil.
// The following fields are editable: Title, Body, State, Base.Ref and MaintainerCanModify.
// Base.Ref updates the base branch of the pull request.
func PullsEdit(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner               string `json:"owner"  zeit:"required"`
			Repo                string `json:"repo"   zeit:"required"`
			Number              int    `json:"number" zeit:"required"`
			Title               string `json:"title"`
			Body                string `josn:"body"`
			State               string `json:"state"`
			Base                string `json:"base"`
			MaintainerCanModify bool   `json:"maintainer_can_modify"`
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
t := c.Request.Header.Get("Authorization")
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

	base := &gogithub.PullRequestBranch{
		Ref: &github.Base,
	}

	opt := &gogithub.PullRequest{
		Title:               &github.Title,
		Body:                &github.Body,
		Base:                base,
		MaintainerCanModify: &github.MaintainerCanModify,
		State:               &github.State,
	}

	pull, _, err := client.Client.PullRequests.Edit(ctx, github.Owner, github.Repo, github.Number, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request": pull})
}
