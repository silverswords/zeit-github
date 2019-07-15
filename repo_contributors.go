package github

import (
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
)

// Contributor API for a list of Contributors in the repositiry
func Contributor(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner string `json:"owner" zeit:"required"`
			Repo  string `json:"repo"  zeit:"required"`
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

	Token := c.Request.Header
	s := cloudapi.NewService("token " + Token["Authorization"][0])

	contributor, err := s.Contributor(github.Owner, github.Repo)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "contributors": contributor})
}
