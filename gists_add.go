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

// GistAdd creates a gist for a user.
func GistAdd(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Public      bool                  `json:"public"      zeit:"required"`
			Description string                `json:"description" zeit:"required"`
			FileName    gogithub.GistFilename `json:"filename"    zeit:"required"`
			Content     string                `json:"cotent"      zeit:"required"`
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

	file := map[gogithub.GistFilename]gogithub.GistFile{
		github.FileName: gogithub.GistFile{
			Content: &github.Content,
		},
	}

	input := &gogithub.Gist{
		Public:      &github.Public,
		Description: &github.Description,
		Files:       file,
	}

	gist, _, err := client.Client.Gists.Create(ctx, input)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "gist": gist})
}