package github

import (
	"context"
	"net/http"
	"time"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// CommitsList lists the commits of a repository
func CommitsList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner   string    `json:"owner"`
			Repo    string    `json:"repo"  zeit:"required"`
			SHA     string    `json:"sha"`
			Path    string    `json:"path"`
			Author  string    `json:"author"`
			Since   time.Time `json:"since"`
			Until   time.Time `josn:"until"`
			Page    int       `json:"page"`
			PerPage int       `josn:"per_page"`
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

	options := gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	opt := &gogithub.CommitsListOptions{
		SHA:         github.SHA,
		Path:        github.Path,
		Author:      github.Author,
		Since:       github.Since,
		Until:       github.Until,
		ListOptions: options,
	}

	commits, _, err := client.Client.Repositories.ListCommits(ctx, github.Owner, github.Repo, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "commits": commits})
}
