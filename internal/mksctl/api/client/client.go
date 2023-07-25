package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	api "github.com/JulienBreux/mksctl/internal/mksctl/api/gen"
)

const (
	defaultAPITimeout = 5 * time.Second

	apiPath = "/api"
)

// Client represents an auth client interface
type Client interface {
	Actions() api.ClientWithResponsesInterface
	APIURL() string
	AuthConfig() (*api.KeycloakConfig, error)
}

// New returns a new instance of Client interface
func New(apiURL string) (Client, error) {
	cwr, err := api.NewClientWithResponses(
		cleanAPIURL(apiURL),
		api.WithHTTPClient(httpClient(defaultAPITimeout)),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		actions: cwr,
		apiURL:  cleanAPIURL(apiURL),
	}, nil
}

type client struct {
	actions *api.ClientWithResponses

	apiURL string
}

// Actions returns possible actions in server API
func (c *client) Actions() api.ClientWithResponsesInterface {
	return c.actions
}

// APIURL returns API URL
func (c *client) APIURL() string {
	return c.apiURL
}

// AuthConfig returns authentication configuration
func (c *client) AuthConfig() (*api.KeycloakConfig, error) {
	ctx := context.Background()
	resp, err := c.actions.GetKeycloakConfigWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func httpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func cleanAPIURL(apiURL string) string {
	apiURL = strings.TrimRight(apiURL, apiPath)
	apiURL = fmt.Sprintf("%s%s", apiURL, apiPath)
	return apiURL
}
