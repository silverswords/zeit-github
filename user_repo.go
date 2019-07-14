package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	con "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// RepoClient encapsulate github.Client
type RepoClient struct {
	GitHubClient *gg.Client
}

// NewRepoClient create RepoClient
func NewRepoClient(g *http.Client) *RepoClient {
	client := gg.NewClient(g)
	return &RepoClient{
		GitHubClient: client,
	}
}

// List list the repositories for a user.
func List(w http.ResponseWriter, r *http.Request) {
	c := con.NewContext(w, r)

	Token := c.Request.Header
	s := Token["Authorization"][0]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := NewRepoClient(tc)

	opt := gg.RepositoryListOptions{
		Visibility:  "all",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
	}

	repolist, _, err := client.GitHubClient.Repositories.List(ctx, "", &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "repolist": repolist})
}
