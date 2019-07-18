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

// PullsMerge merges a pull request (Merge Button™).
// commitMessage is the title for the automatic commit message.
func PullsMerge(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner         string `json:"owner"          zeit:"required"`
			Repo          string `json:"repo"           zeit:"required"`
			Number        int    `json:"number"         zeit:"required"`
			CommitMessage string `json:"commit_message" zeit:"required"`
			CommitTitle   string `json:"commit_title"`
			SHA           string `json:"sha"`
			MergeMethod   string `json:"merge_method"`
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

	opt := &gogithub.PullRequestOptions{
		CommitTitle: github.CommitTitle,
		SHA:         github.SHA,
		MergeMethod: github.MergeMethod,
	}

	pull, _, err := client.Client.PullRequests.Merge(ctx, github.Owner, github.Repo, github.Number, github.CommitMessage, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request": pull})
}
