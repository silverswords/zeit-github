package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// GistClient -
type GistClient struct {
	GitHubClient *gg.Client
}

// NewGistClient create GithubCliebt
func NewGistClient(g *http.Client) *GistClient {
	client := gg.NewClient(g)
	return &GistClient{
		GitHubClient: client,
	}
}

// GistList -
func GistList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			User string `json:"user" zeit:"required"`
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

	client := NewGistClient(nil)
	ctx := context.Background()

	gist, _, err := client.GitHubClient.Gists.List(ctx, github.User, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "Gists": gist})
}
