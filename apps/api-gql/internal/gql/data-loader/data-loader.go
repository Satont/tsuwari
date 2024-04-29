package data_loader

import (
	"context"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type DataLoader struct {
	cachedTwitchClient    *twitch.CachedTwitchClient
	helixUserByIdLoader   *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
	helixUserByNameLoader *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
}

type Opts struct {
	CachedTwitchClient *twitch.CachedTwitchClient
}

func New(opts Opts) *DataLoader {
	loader := &DataLoader{
		cachedTwitchClient: opts.CachedTwitchClient,
	}

	loader.helixUserByIdLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByIds,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.helixUserByNameLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByNames,
		dataloadgen.WithWait(time.Millisecond),
	)

	return loader
}

// For returns the dataloader for a given context
func GetLoaderForRequest(ctx context.Context) *DataLoader {
	return ctx.Value(LoadersKey).(*DataLoader)
}
