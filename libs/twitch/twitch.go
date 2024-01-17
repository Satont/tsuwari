package twitch

import (
	"context"
	"time"

	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/nicklaw5/helix/v2"
)

type helixClient struct {
	*helix.Client
}

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	var reset64 int64
	reset64 = int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().UTC().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}

func NewAppClient(config cfg.Config, tokensGrpc tokens.TokensClient) (*helix.Client, error) {
	return NewAppClientWithContext(context.Background(), config, tokensGrpc)
}

func NewAppClientWithContext(
	ctx context.Context,
	config cfg.Config,
	tokensGrpc tokens.TokensClient,
) (
	*helix.Client, error,
) {
	appToken, err := tokensGrpc.RequestAppToken(
		ctx,
		&emptypb.Empty{},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:       config.TwitchClientId,
			ClientSecret:   config.TwitchClientSecret,
			RedirectURI:    config.TwitchCallbackUrl,
			RateLimitFunc:  rateLimitCallback,
			AppAccessToken: appToken.AccessToken,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewUserClient(userID string, config cfg.Config, tokensGrpc tokens.TokensClient) (
	*helix.Client,
	error,
) {
	return NewUserClientWithContext(context.Background(), userID, config, tokensGrpc)
}

func NewUserClientWithContext(
	ctx context.Context,
	userID string,
	config cfg.Config,
	tokensGrpc tokens.TokensClient,
) (*helix.Client, error) {
	userToken, err := tokensGrpc.RequestUserToken(
		ctx,
		&tokens.GetUserTokenRequest{UserId: userID},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:        config.TwitchClientId,
			ClientSecret:    config.TwitchClientSecret,
			RedirectURI:     config.TwitchCallbackUrl,
			RateLimitFunc:   rateLimitCallback,
			UserAccessToken: userToken.AccessToken,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewBotClient(botID string, config cfg.Config, tokensGrpc tokens.TokensClient) (
	*helix.Client,
	error,
) {
	return NewBotClientWithContext(context.Background(), botID, config, tokensGrpc)
}

func NewBotClientWithContext(
	ctx context.Context, botID string, config cfg.Config, tokensGrpc tokens.TokensClient,
) (*helix.Client, error) {
	botToken, err := tokensGrpc.RequestBotToken(
		ctx,
		&tokens.GetBotTokenRequest{BotId: botID},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:        config.TwitchClientId,
			ClientSecret:    config.TwitchClientSecret,
			RedirectURI:     config.TwitchCallbackUrl,
			RateLimitFunc:   rateLimitCallback,
			UserAccessToken: botToken.AccessToken,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
