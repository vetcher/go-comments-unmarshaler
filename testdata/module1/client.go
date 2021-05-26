package module1

import (
	"context"
)

// And this Client is from module1
type Client struct{}

// Comment for another Do
func (Client) Do(ctx context.Context) error { panic("not implemented") }
