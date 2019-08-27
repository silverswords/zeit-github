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

// RepoProjectAdd creates a GitHub Project for the specified repository.
func RepoProjectAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner                  string `json:"owner"  zeit:"required"`
			Repo                   string `json:"repo"   zeit:"required"`
			Name                   string `json:"name"   zeit:"required"`
			Body                   string `json:"body"`
			State                  string `json:"state"`
			Public                 bool   `json:"pulic"`
			OrganizationPermission string `json:"organization_permission"`
		}

		tc *http.Client
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

	ctx := context.Background()

	t := c.Request.Header.Get("Authorization")
	if t != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: t},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := cloudapi.NewAPIClient(tc)

	opt := &gogithub.ProjectOptions{
		Name:                   &github.Name,
		Body:                   &github.Body,
		State:                  &github.State,
		OrganizationPermission: &github.OrganizationPermission,
		Public:                 &github.Public,
	}

	contributor, _, err := client.Client.Repositories.CreateProject(ctx, github.Owner, github.Repo, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "contributors": contributor})
}
