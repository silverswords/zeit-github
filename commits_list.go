package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// ComClient -
type ComClient struct {
	GitHubClient *gg.Client
}

// NewClient create GithubCliebt
func NewClient(g *http.Client) *ConClient {
	client := gg.NewClient(g)
	return &ConClient{
		GitHubClient: client,
	}
}

// CommitsList -
func CommitsList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner string
			Repo  string `zeit:"required"`
		}
	)
	c := con.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusNotAcceptable, con.H{"status": http.StatusNotAcceptable})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusConflict, con.H{"status": http.StatusConflict})
		return
	}

	Token := c.Request.Header
	s := Token["Authorization"][0]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := NewListClient(tc)

	commits, _, err := client.GitHubClient.Repositories.ListCommits(ctx, github.Owner, github.Repo, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "Commits": commits})
}
