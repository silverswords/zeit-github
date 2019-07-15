package github

import (
	"context"
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkg "github.com/silverswords/clouds/pkgs/http/context"
)

// Following lists the people that a user is following.
func Following(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			User string `json:"user" zeit:"required"`
		}
	)

	c := cloudpkg.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusNotAcceptable, cloudpkg.H{"status": http.StatusNotAcceptable})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusConflict, cloudpkg.H{"status": http.StatusConflict})
		return
	}

	client := cloudapi.NewAPIClient(nil)
	ctx := context.Background()

	following, _, err := client.Client.Users.ListFollowers(ctx, github.User, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkg.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkg.H{"status": http.StatusOK, "following": following})
}
