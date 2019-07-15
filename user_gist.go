package github

import (
	"context"
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkg "github.com/silverswords/clouds/pkgs/http/context"
)

// GistList  list gists for a user.
func GistList(w http.ResponseWriter, r *http.Request) {
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

	gist, _, err := client.Client.Gists.List(ctx, github.User, nil)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkg.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkg.H{"status": http.StatusOK, "gists": gist})
}
