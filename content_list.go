package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// ConClient -
type ConClient struct {
	GitHubClient *gg.Client
}

// NewConClient create GithubCliebt
func NewConClient(g *http.Client) *ConClient {
	client := gg.NewClient(g)
	return &ConClient{
		GitHubClient: client,
	}
}

// ContentList -
func ContentList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner string `json:"owner"`
			Repo  string `json:"repo" zeit:"required"`
			Path  string `json:"path"`
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
	client := NewConClient(tc)

	_, contentList, _, err := client.GitHubClient.Repositories.GetContents(ctx, github.Owner, github.Repo, github.Path, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "Contents": contentList})
}
