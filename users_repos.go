package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
)

// UsersReposList lists the repositories for a user.
func UsersReposList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			UserName string `json:"user_name" zeit:"required"`
			Type     string `json:"type"`
			Sort     string `json:"sort"`
			Page     int    `json:"page"`
			PerPage  int    `josn:"per_page"`
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

	client := cloudapi.NewAPIClient(nil)
	ctx := context.Background()

	options := gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	opt := gogithub.RepositoryListOptions{
		Type:        github.Type,
		Sort:        github.Sort,
		ListOptions: options,
	}

	repolist, _, err := client.Client.Repositories.List(ctx, github.UserName, &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": err})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "repo_list": repolist})
}
