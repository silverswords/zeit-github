package github

import (
	"context"
	"net/http"
	"time"

	gogithub "github.com/google/go-github/v27/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// IssuesList lists the issues for the authenticated user. If all is true, list issues
// across all the user's visible repositories including owned, member, and
// organization repositories; if false, list only owned and member
// repositories.
func IssuesList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			All       bool      `json:"all"     zeit:"required"`
			Filter    string    `json:"filter"`
			State     string    `json:"state"`
			Labels    []string  `json:"labels"`
			Sort      string    `json:"sort"`
			Direction string    `json:"direction"`
			Since     time.Time `json:"since"`
			Page      int       `json:"page"`
			PerPage   int       `json:"per_page"`
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

	options := gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	opt := &gogithub.IssueListOptions{
		Filter:      github.Filter,
		State:       github.State,
		Labels:      github.Labels,
		Sort:        github.Sort,
		Direction:   github.Direction,
		Since:       github.Since,
		ListOptions: options,
	}

	issueList, _, err := client.Client.Issues.List(ctx, github.All, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "issues": issueList})
}
