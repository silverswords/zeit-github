package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// ListStarred lists all the repos starred by a user. Passing the empty string
// will list the starred repositories for the authenticated user.
func ListStarred(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			User      string `json:"user"`
			Sort      string `json:"sort"`
			Direction string `json:"direction"`
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

	token := c.Request.Header
	t := token.Get("Authorization")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	opt := gogithub.ActivityListStarredOptions{
		Sort:      github.Sort,
		Direction: github.Direction,
	}

	star, _, err := client.Client.Activity.ListStarred(ctx, github.User, &opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "star": star})
}
