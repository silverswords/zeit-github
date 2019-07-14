package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// FollowerClient encapsulate github.Client
type FollowerClient struct {
	GitHubClient *gg.Client
}

// NewFollowerClient create FollowerClient
func NewFollowerClient(g *http.Client) *FollowerClient {
	client := gg.NewClient(g)
	return &FollowerClient{
		GitHubClient: client,
	}
}

// Follower lists the followers for a user.
func Follower(w http.ResponseWriter, r *http.Request) {
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

	client := NewFollowerClient(nil)
	ctx := context.Background()

	followers, _, err := client.GitHubClient.Users.ListFollowers(ctx, github.User, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "followers": followers})
}
