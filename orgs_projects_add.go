package github

import (
	"context"
	"net/http"

	gogithub "github.com/google/go-github/github"
	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// OrgsProjectsAdd lists the projects for an organization.
func OrgsProjectsAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Org                    string `json:"org"   zeit:"required"`
			Name                   string `json:"name"   zeit:"required"`
			Body                   string `json:"body"   zeit:"required"`
			State                  string `json:"state"`
			Public                 bool   `json:"pulic"`
			OrganizationPermission string `json:"organizationPermission"`
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

	opt := &gogithub.ProjectOptions{
		Name:                   &github.Name,
		Body:                   &github.Body,
		State:                  &github.State,
		OrganizationPermission: &github.OrganizationPermission,
		Public:                 &github.Public,
	}
	project, _, err := client.Client.Organizations.CreateProject(ctx, github.Org, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "project": project})
}
