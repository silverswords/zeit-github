package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// SeaClient -
type SeaClient struct {
	GitHubClient *gg.Client
}

// NewSeaClient create GithubCliebt
func NewSeaClient(g *http.Client) *SeaClient {
	client := gg.NewClient(g)
	return &SeaClient{
		GitHubClient: client,
	}
}

// Search -
func Search(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Key  string `json:"key" zeit:"required"`
			Sort string `json:"sort" zeit:"required"`
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

	client := NewSeaClient(nil)
	ctx := context.Background()

	Options := gg.ListOptions{Page: 1, PerPage: 10}
	opts := &gg.SearchOptions{Sort: github.Sort, Order: "desc", ListOptions: Options}

	repo, _, err := client.GitHubClient.Search.Repositories(ctx, github.Key, opts)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "a": repo})
}
