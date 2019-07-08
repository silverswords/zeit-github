package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	con "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// Client -
type Client struct {
	GitHubClient *gg.Client
}

// NewClient create GithubCliebt
func NewClient(g *http.Client) *Client {
	client := gg.NewClient(g)
	return &Client{
		GitHubClient: client,
	}
}

// List -
func List(w http.ResponseWriter, r *http.Request) {
	c := con.NewContext(w, r)

	Token := c.Request.Header
	s := Token["Authorization"][0]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := NewClient(tc)

	opt := gg.RepositoryListOptions{
		Visibility:  "all",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
	}

	list, _, err := client.GitHubClient.Repositories.List(ctx, "", &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "list": list})
}
