package current

import (
	"context"

	"github.com/zmb3/spotify"
)

var client *spotify.Client

func SpotifyClient(ctx context.Context) *spotify.Client {
	v := ctx.Value(clientCtxKey)
	if v != nil {
		return v.(*spotify.Client)
	}

	return client
}

func WithSpotifyClient(ctx context.Context, client *spotify.Client) context.Context {
	return context.WithValue(ctx, clientCtxKey, client)
}

func SetSpotifyClient(cl *spotify.Client) {
	if client != nil {
		panic("cannot call SetClient twice")
	}
	if cl == nil {
		panic("cl must not be nil")
	}
	client = cl
}
