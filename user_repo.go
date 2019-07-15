package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// List list the repositories for a user.
func List(w http.ResponseWriter, r *http.Request) {
	c := cloudpkgs.NewContext(w, r)

	token := c.Request.Header
	t := token.Get("Authorization")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	opt := gogithub.RepositoryListOptions{
		Visibility:  "all",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
	}

	repolist, _, err := client.Client.Repositories.List(ctx, "", &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "repolist": repolist})
}
