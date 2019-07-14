package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// SearchRepoClient encapsulate github.Client
type SearchRepoClient struct {
	GitHubClient *gg.Client
}

// NewSearchRepoClient create SearchRepoClient
func NewSearchRepoClient(g *http.Client) *SearchRepoClient {
	client := gg.NewClient(g)
	return &SearchRepoClient{
		GitHubClient: client,
	}
}

// SearchRepo searches repositories via various criteria.
func SearchRepo(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Key  string `json:"key"  zeit:"required"`
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

	client := NewSearchRepoClient(nil)
	ctx := context.Background()

	Options := gg.ListOptions{Page: 1, PerPage: 20}
	opts := &gg.SearchOptions{Sort: github.Sort, Order: "desc", ListOptions: Options}

	repo, _, err := client.GitHubClient.Search.Repositories(ctx, github.Key, opts)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "repo": repo})
}
