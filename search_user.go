package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// SearchUserClient encapsulate github.Client
type SearchUserClient struct {
	GitHubClient *gg.Client
}

// NewSearchUserClient create GithubClient
func NewSearchUserClient(g *http.Client) *SearchUserClient {
	client := gg.NewClient(g)
	return &SearchUserClient{
		GitHubClient: client,
	}
}

// SearchUser searches users via various criteria.
func SearchUser(w http.ResponseWriter, r *http.Request) {
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

	client := NewSearchUserClient(nil)
	ctx := context.Background()

	Options := gg.ListOptions{Page: 1, PerPage: 20}
	opts := &gg.SearchOptions{Sort: github.Sort, Order: "desc", ListOptions: Options}

	user, _, err := client.GitHubClient.Search.Users(ctx, github.Key, opts)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "user": user})
}
