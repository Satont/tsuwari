package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// IntegrationsGetServiceAuthLink is the resolver for the integrationsGetServiceAuthLink field.
func (r *queryResolver) IntegrationsGetServiceAuthLink(
	ctx context.Context,
	service gqlmodel.IntegrationService,
) (string, error) {
	return "", nil
}

// IntegrationsGetData is the resolver for the integrationsGetData field.
func (r *queryResolver) IntegrationsGetData(
	ctx context.Context,
	service gqlmodel.IntegrationService,
) (gqlmodel.IntegrationData, error) {
	switch service {
	case gqlmodel.IntegrationServiceLastfm:
		return gqlmodel.IntegrationDataLastfm{
			Username: "lastfm",
			Avatar:   "123",
		}, nil
	case gqlmodel.IntegrationServiceVk:
		return gqlmodel.IntegrationDataVk{
			Username: "vk",
			Avatar:   "321",
		}, nil
	case gqlmodel.IntegrationServiceSpotify:
		return gqlmodel.IntegrationDataSpotify{
			Username: "spotify",
			Avatar:   "1",
		}, nil
	case gqlmodel.IntegrationServiceDonationalerts:
		return gqlmodel.IntegrationDataDonationAlerts{
			Username: "donationalerts",
			Avatar:   "1",
		}, nil
	case gqlmodel.IntegrationServiceDiscord:
		return gqlmodel.IntegrationDataDiscord{
			Guilds: []gqlmodel.IntegrationDataDiscordGuild{},
		}, nil
	case gqlmodel.IntegrationServiceStreamlabs:
		return gqlmodel.IntegrationDataStreamLabs{
			Username: "streamlabs",
			Avatar:   "",
		}, nil
	case gqlmodel.IntegrationServiceValorant:
		return gqlmodel.IntegrationDataValorant{
			Username: "valorant",
			Avatar:   "1",
		}, nil
	}

	return nil, fmt.Errorf("unknown service: %s", service)
}
