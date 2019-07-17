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

// IssueEdit  edits an issue.
func IssueEdit(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner     string   `json:"owner"  zeit:"required"`
			Repo      string   `json:"repo"   zeit:"required"`
			Number    int      `json:"number" zeit:"required"`
			Title     string   `json:"title"`
			Body      string   `json:"body"`
			Labels    []string `json:"labels"`
			Assignee  string   `json:"assignee"`
			State     string   `json:"state"`
			Milestone int      `json:"milestone"`
			Assignees []string `json:"assignees"`
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

	req := &gogithub.IssueRequest{
		Title:     &github.Title,
		Body:      &github.Body,
		Labels:    &github.Labels,
		Assignee:  &github.Assignee,
		State:     &github.State,
		Milestone: &github.Milestone,
		Assignees: &github.Assignees,
	}

	issue, _, err := client.Client.Issues.Edit(ctx, github.Owner, github.Repo, github.Number, req)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "issue": issue})
}
