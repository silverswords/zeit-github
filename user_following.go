package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// FollowingClient encapsulate github.Client
type FollowingClient struct {
	GitHubClient *gg.Client
}

// NewFollowingClient create FollowingClient
func NewFollowingClient(g *http.Client) *FollowingClient {
	client := gg.NewClient(g)
	return &FollowingClient{
		GitHubClient: client,
	}
}

// Following lists the people that a user is following.
func Following(w http.ResponseWriter, r *http.Request) {
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

	client := NewFollowingClient(nil)
	ctx := context.Background()

	following, _, err := client.GitHubClient.Users.ListFollowers(ctx, github.User, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "following": following})
}
