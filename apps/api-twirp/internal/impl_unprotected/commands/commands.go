package commands

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/commands_unprotected"
)

type Commands struct {
	*impl_deps.Deps
}

func (c *Commands) GetChannelCommands(ctx context.Context, req *commands_unprotected.GetChannelCommandsRequest) (*commands_unprotected.GetChannelCommandsResponse, error) {
	return nil, nil
}
