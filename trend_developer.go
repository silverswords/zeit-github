package github

import (
	"net/http"

	service "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	con "github.com/silverswords/clouds/pkgs/http/context"
)

// DeveloperTrend return an array of trending developers.
func DeveloperTrend(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Language  string `json:"language" `
			DataRange string `json:"datarange"`
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

	trend, err := service.DeveloperTrending(github.Language, github.DataRange)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, con.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, con.H{"status": http.StatusOK, "trend": trend})

}
