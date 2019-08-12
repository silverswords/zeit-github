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

// OrgsNumbersList lists the members for an organization. If the authenticated
// user is an owner of the organization, this will return both concealed and
// public members, otherwise it will only return public members.
func OrgsNumbersList(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Org        string `json:"org"         zeit:"required"`
			PublicOnly bool   `json:"public_only"`
			Filter     string `json:"filter"`
			Role       string `json:"role"`
			Page       int    `json:"page"`
			PerPage    int    `json:"per_page"`
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

	options := gogithub.ListOptions{
		Page:    github.Page,
		PerPage: github.PerPage,
	}

	opt := &gogithub.ListMembersOptions{
		PublicOnly:  github.PublicOnly,
		Filter:      github.Filter,
		Role:        github.Role,
		ListOptions: options,
	}

	numbers, _, err := client.Client.Organizations.ListMembers(ctx, github.Org, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "numbers": numbers})
}
