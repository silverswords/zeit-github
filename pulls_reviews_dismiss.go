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

// PullsReviewsDismiss dismisses a specified review on the specified pull request.
func PullsReviewsDismiss(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner    string `json:"owner"     zeit:"required"`
			Repo     string `json:"repo"      zeit:"required"`
			Number   int    `json:"number"    zeit:"required"`
			ReviewID int64  `json:"review_id" zeit:"required"`
			Message  string `json:"message"`
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

	t := c.Request.Header.Get("Authorization")
	if t == "" {
		c.WriteJSON(http.StatusUnauthorized, cloudpkgs.H{"status": http.StatusUnauthorized})
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	opt := &gogithub.PullRequestReviewDismissalRequest{
		Message: &github.Message,
	}

	pull, _, err := client.Client.PullRequests.DismissReview(ctx, github.Owner, github.Repo, github.Number, github.ReviewID, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request_review": pull})
}
