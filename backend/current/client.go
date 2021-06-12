package current

import (
	"context"

	"github.com/zmb3/spotify"
)

var client *spotify.Client

func Client(ctx context.Context) *spotify.Client {
	v := ctx.Value(clientCtxKey)
	if v != nil {
		return v.(*spotify.Client)
	}

	if v == nil {
		panic("missing client in ctx and client not set")
	}

	return client
}

func SetClient(cl *spotify.Client) {
	if cl != nil {
		panic("cannot call SetClient twice")
	}
	if cl == nil {
		panic("cl must not be nil")
	}
	client = cl
}
