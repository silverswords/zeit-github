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

// UsersSelfRepos lists repositories for the authenticated user.
func UsersSelfRepos(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Visibility  string `json:"visibility"`
			Affiliation string `json:"affiliation"`
			Direction   string `json:"direction"`
			Type        string `json:"type"`
			Sort        string `json:"sort"`
			Page        int    `json:"page"`
			PerPage     int    `josn:"per_page"`
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

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	options := gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	opt := gogithub.RepositoryListOptions{
		Visibility:  github.Visibility,
		Affiliation: github.Affiliation,
		Type:        github.Type,
		Sort:        github.Sort,
		Direction:   github.Direction,
		ListOptions: options,
	}

	repolist, _, err := client.Client.Repositories.List(ctx, "", &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "repo_list": repolist})
}
