package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
)

// SearchRepo searches repositories via various criteria.
func SearchRepo(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Key     string `json:"key"      zeit:"required"`
			Sort    string `json:"sort"     zeit:"required"`
			Page    int    `json:"page"     zeit:"required"`
			PerPage int    `json:"per_page" zeit:"required"`
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

	Options := gogithub.ListOptions{Page: github.Page, PerPage: github.PerPage}
	opts := &gogithub.SearchOptions{Sort: github.Sort, Order: "desc", ListOptions: Options}

	repo, _, err := client.Client.Search.Repositories(ctx, github.Key, opts)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "repo": repo})
}
