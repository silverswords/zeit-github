package github

import (
	"context"
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// RateLimit returns the rate limits for the current client.
func RateLimit(w http.ResponseWriter, r *http.Request) {
	c := cloudpkgs.NewContext(w, r)

	token := c.Request.Header
	t := token.Get("Authorization")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	if t == "" {
		tc = nil
	}
	client := cloudapi.NewAPIClient(tc)

	rate, _, err := client.Client.RateLimits(ctx)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "rate": rate})
}
