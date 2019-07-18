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

// ReleasesEdit edits a repository release.
// Note that only a subset of the release fields are used.
// See RepositoryRelease for more information.
func ReleasesEdit(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner           string `json:"owner"  zeit:"required"`
			Repo            string `json:"repo"   zeit:"required"`
			ID              int64  `json:"id"     zeit:"required"`
			TagName         string `json:"tag_name"`
			TargetCommitish string `json:"target_commitish"`
			Name            string `json:"name"`
			Body            string `josn:"body"`
			Draft           bool   `json:"draft"`
			Prerelease      bool   `json:"prerelease"`
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

	opt := &gogithub.RepositoryRelease{
		TagName:         &github.TagName,
		Name:            &github.Name,
		Body:            &github.Body,
		TargetCommitish: &github.TargetCommitish,
		Draft:           &github.Draft,
		Prerelease:      &github.Prerelease,
	}

	release, _, err := client.Client.Repositories.EditRelease(ctx, github.Owner, github.Repo, github.ID, opt)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "release": release})
}
