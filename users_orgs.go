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

// UserOrgs lists the organizations for a user. Passing the empty string will list
// organizations for the authenticated user.
func UserOrgs(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			User    string `json:"user"`
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

	if github.User == "" {
		token := c.Request.Header
		t := token.Get("Authorization")

		if t == "" {
			c.WriteJSON(http.StatusUnauthorized, cloudpkgs.H{"status": http.StatusUnauthorized})
			return
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: t},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := cloudapi.NewAPIClient(tc)

	opt := &gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	org, _, err := client.Client.Organizations.List(ctx, github.User, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "orgs": org})
}
