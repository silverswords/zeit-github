package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
)

// UserFollowing lists the followers for a user. Passing the empty string will
// fetch followers for the authenticated user.
func UserFollowing(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			User    string `json:"user"`
			Page    int    `json:"page"`
			PerPage int    `json:"per_page"`
		}
	)

	c := cloudpkgs.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusBadRequest, cloudpkgs.H{"status": http.StatusBadRequest})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusPreconditionRequired, cloudpkgs.H{"status": http.StatusPreconditionRequired})
		return
	}

	client := cloudapi.NewAPIClient(nil)
	ctx := context.Background()

	options := &gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	following, _, err := client.Client.Users.ListFollowers(ctx, github.User, options)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "following": following})
}
