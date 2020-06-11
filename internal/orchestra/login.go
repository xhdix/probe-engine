package orchestra

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/internal/httpx"
)

// LoginCredentials contains the login credentials
type LoginCredentials struct {
	ClientID string `json:"username"`
	Password string `json:"password"`
}

// LoginAuth contains authentication info
type LoginAuth struct {
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}

// MaybeLogin performs login if necessary
func (c Client) MaybeLogin(ctx context.Context) error {
	state := c.StateFile.Get()
	if state.Auth() != nil {
		return nil // we're already good
	}
	creds := state.Credentials()
	if creds == nil {
		return errNotRegistered
	}
	c.LoginCalls.Add(1)
	var auth LoginAuth
	err := (httpx.Client{
		BaseURL:    c.BaseURL,
		HTTPClient: c.HTTPClient,
		Logger:     c.Logger,
		UserAgent:  c.UserAgent,
	}).CreateJSON(ctx, "/api/v1/login", *creds, &auth)
	if err != nil {
		return err
	}
	state.Expire = auth.Expire
	state.Token = auth.Token
	return c.StateFile.Set(state)
}