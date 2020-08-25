package provider

import (
	"errors"

	"github.com/rainproj/rain/internal/providers/dockerhub"
	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
)

// Provider interface
type Provider interface {
	Push(*context.Context, config.Push) error
}

var providers = map[string]Provider{
	"hub": dockerhub.Provider{},
}

// Get provider
func Get(p string) (Provider, error) {
	pusher, ok := providers[p]
	if !ok {
		return nil, errors.New("Unknown provider")
	}

	return pusher, nil
}
