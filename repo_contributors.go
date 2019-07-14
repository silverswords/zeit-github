package github

import (
	"net/http"

	service "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// Contributor API for a list of Contributors in the repositiry
func Contributor(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner string `json:"owner" zeit:"required"`
			Repo  string `json:"repo" zeit:"required"`
		}
	)

	c := con.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusNotAcceptable, con.H{"status": http.StatusNotAcceptable})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusConflict, con.H{"status": http.StatusConflict})
		return
	}

	Token := c.Request.Header
	s := service.NewService("token " + Token["Authorization"][0])

	contributor, err := s.Contributor(github.Owner, github.Repo)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": err})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "contributors": contributor})
}
