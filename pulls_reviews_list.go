package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// PullsReviewsList lists all reviews on the specified pull request.
func PullsReviewsList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner   string `json:"owner"   zeit:"required"`
			Repo    string `json:"repo"    zeit:"required"`
			Number  int    `json:"number"  zeit:"required"`
			Page    int    `json:"page"`
			PerPage int    `json:"per_page"`
		}

		tc *http.Client
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

	ctx := context.Background()

	t := c.Request.Header.Get("Authorization")
	if t != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: t},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := cloudapi.NewAPIClient(tc)

	options := &gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	pull, _, err := client.Client.PullRequests.ListReviews(ctx, github.Owner, github.Repo, github.Number, options)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request_review": pull})
}
