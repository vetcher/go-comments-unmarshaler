package client

import (
	"context"
)

// Client is client from package `client`
type Client struct{}

// Do some stuff
func (Client) Do(ctx context.Context) error { panic("not implemented") }
