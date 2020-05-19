// Package probeservices contains code to contact OONI probe services.
//
// Specifically we implement v2.0.0 of the OONI bouncer specification defined
// in https://github.com/ooni/spec/blob/master/backends/bk-004-bouncer
//
// We additionally implement v2.0.0 of the OONI collector specification defined
// in https://github.com/ooni/spec/blob/master/backends/bk-003-collector.md.
package probeservices

import (
	"context"

	"github.com/ooni/probe-engine/internal/jsonapi"
	"github.com/ooni/probe-engine/model"
)

// Client is a client for the OONI probe services API.
type Client struct {
	jsonapi.Client
}

// GetCollectors queries the bouncer for collectors. Returns a list of
// entries on success; an error on failure.
func (c *Client) GetCollectors(ctx context.Context) (output []model.Service, err error) {
	err = c.Client.Read(ctx, "/api/v1/collectors", &output)
	return
}

// GetTestHelpers is like GetCollectors but for test helpers.
func (c *Client) GetTestHelpers(
	ctx context.Context) (output map[string][]model.Service, err error) {
	err = c.Client.Read(ctx, "/api/v1/test-helpers", &output)
	return
}
