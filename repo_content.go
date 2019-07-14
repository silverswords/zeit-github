package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// ContentClient encapsulate github.Client
type ContentClient struct {
	GitHubClient *gg.Client
}

// NewContentClient create ContentClient
func NewContentClient(g *http.Client) *ContentClient {
	client := gg.NewClient(g)
	return &ContentClient{
		GitHubClient: client,
	}
}

// ContentList return either the metadata and content of a single file
// (when path references a file) or the metadata of all the files and/or
// subdirectories of a directory (when path references a directory)
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
	client := NewContentClient(tc)

	_, contentList, _, err := client.GitHubClient.Repositories.GetContents(ctx, github.Owner, github.Repo, github.Path, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "Contents": contentList})
}
