package github

import (
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
)

// RepoTrend return an array of trending repositories.
func RepoTrend(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Language  string `json:"language"`
			DataRange string `json:"daterange"`
		}
	)

	c := cloudpkgs.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusNotAcceptable, cloudpkgs.H{"status": http.StatusNotAcceptable})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusConflict, cloudpkgs.H{"status": http.StatusConflict})
		return
	}

	trend, err := cloudapi.RepoTrending(github.Language, github.DataRange)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "trend": trend})

}
