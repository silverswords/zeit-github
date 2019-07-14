package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// SearchCodeClient encapsulate github.Client
type SearchCodeClient struct {
	GitHubClient *gg.Client
}

// NewSearchCodeClient create SearchCodeClient
func NewSearchCodeClient(g *http.Client) *SearchCodeClient {
	client := gg.NewClient(g)
	return &SearchCodeClient{
		GitHubClient: client,
	}
}

// SearchCode searches code via various criteria.
func SearchCode(w http.ResponseWriter, r *http.Request) {
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

	client := NewSearchCodeClient(nil)
	ctx := context.Background()

	Options := gg.ListOptions{Page: 1, PerPage: 50}
	opts := &gg.SearchOptions{Sort: github.Sort, Order: "desc", ListOptions: Options}

	code, _, err := client.GitHubClient.Search.Code(ctx, github.Key, opts)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "code": code})
}
