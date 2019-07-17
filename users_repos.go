package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// List list the repositories for a user.
func List(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Visibility  string `json:"visibility"`
			Affiliation string `json:"affiliation"`
			Sort        string `json:"sort"`
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

	opt := gogithub.RepositoryListOptions{
		Visibility:  github.Visibility,
		Affiliation: github.Affiliation,
		Sort:        github.Sort,
		Direction:   "asc",
	}

	repolist, _, err := client.Client.Repositories.List(ctx, "", &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "repo_list": repolist})
}
