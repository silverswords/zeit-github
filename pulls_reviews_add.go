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

// PullsReviewsAdd creates a new review on the specified pull request.
func PullsReviewsAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner    string                         `json:"owner"   zeit:"required"`
			Repo     string                         `json:"repo"    zeit:"required"`
			Number   int                            `json:"number"  zeit:"required"`
			Body     string                         `json:"body"    zeit:"required"`
			CommitID string                         `json:"commit_id"`
			Event    string                         `json:"event"`
			Comments []*gogithub.DraftReviewComment `json:"comments"`
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

	opt := &gogithub.PullRequestReviewRequest{
		Body:     &github.Body,
		CommitID: &github.CommitID,
		Event:    &github.Event,
		Comments: github.Comments,
	}

	pull, _, err := client.Client.PullRequests.CreateReview(ctx, github.Owner, github.Repo, github.Number, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request_review": pull})
}
