package dns

import (
	"context"

	"github.com/qdm12/dns/v2/pkg/cache"
	"github.com/qdm12/dns/v2/pkg/cache/lru"
	"github.com/qdm12/dns/v2/pkg/dot"
	"github.com/qdm12/dns/v2/pkg/filter/mapfilter"
	"github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) GetSettings() (settings settings.DNS) { return l.state.GetSettings() }

func (l *Loop) SetSettings(ctx context.Context, settings settings.DNS) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}

func buildDoTSettings(settings settings.DNS,
	filter *mapfilter.Filter, warner Warner) (
	dotSettings dot.ServerSettings) {
	var cache cache.Interface
	if *settings.DoT.Unbound.Caching {
		cache = lru.New(lru.Settings{})
	}
	providers := make([]provider.Provider, len(settings.DoT.Unbound.Providers))
	for i := range settings.DoT.Unbound.Providers {
		var err error
		providers[i], err = provider.Parse(settings.DoT.Unbound.Providers[i])
		if err != nil {
			panic(err) // this should already been checked
		}
	}

	return dot.ServerSettings{
		Resolver: dot.ResolverSettings{
			DoTProviders: settings.DoT.Unbound.Providers,
			DNSProviders: settings.DoT.Unbound.Providers,
			IPv6:         *settings.DoT.Unbound.IPv6,
			Warner:       warner,
		},
		Filter: filter,
		Cache:  cache,
	}
}
