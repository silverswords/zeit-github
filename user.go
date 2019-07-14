package github

import (
	"context"
	"net/http"

	gg "github.com/google/go-github/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// UserClient -
type UserClient struct {
	GitHubClient *gg.Client
}

// NewUserClient create GithubCliebt
func NewUserClient(g *http.Client) *UserClient {
	client := gg.NewClient(g)
	return &UserClient{
		GitHubClient: client,
	}
}

// User -
func User(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			UserName string `json:"username" zeit:"required"`
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

	client := NewUserClient(nil)
	ctx := context.Background()

	user, _, err := client.GitHubClient.Users.Get(ctx, github.UserName)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "user": user})
}
