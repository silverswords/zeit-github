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

// OrgsEdit  edits an organization.
func OrgsEdit(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			OrgName               string `json:"org_name" zeit:"required"`
			Name                  string `json:"name"`
			Company               string `json:"company"`
			BillingEmail          string `json:"billing_email"`
			Email                 string `json:"email"`
			Location              string `json:"locatin"`
			Description           string `json:"description"`
			DefaultRepoPermission string `json:"default_repo_permission"`
			MembersCanCreateRepos bool   `json:"members_can_create_repositories"`
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

	opt := &gogithub.Organization{
		Name:                  &github.Name,
		Description:           &github.Description,
		Company:               &github.Company,
		BillingEmail:          &github.BillingEmail,
		Email:                 &github.Email,
		Location:              &github.Location,
		DefaultRepoPermission: &github.DefaultRepoPermission,
		MembersCanCreateRepos: &github.MembersCanCreateRepos,
	}

	org, _, err := client.Client.Organizations.Edit(ctx, github.OrgName, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "org": org})
}
