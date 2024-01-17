package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/nicklaw5/helix/v2"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var appTokenScopes = []string{}

type appToken struct {
	AccessToken    string
	ObtainmentTime time.Time
	ExpiresIn      int
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config  cfg.Config
	Gorm    *gorm.DB
	Redsync *redsync.Redsync
	Logger  logger.Logger
}

type tokensGrpcImpl struct {
	tokens.UnimplementedTokensServer

	globalClient   *helix.Client
	appAccessToken *appToken

	appLock   *redsync.Mutex
	usersLock *redsync.Mutex
	botsLock  *redsync.Mutex

	db     *gorm.DB
	config cfg.Config
	log    logger.Logger
}

func NewTokens(opts Opts) error {
	helixClient, err := helix.NewClient(
		&helix.Options{
			ClientID:     opts.Config.TwitchClientId,
			ClientSecret: opts.Config.TwitchClientSecret,
			RedirectURI:  opts.Config.TwitchCallbackUrl,
		},
	)
	if err != nil {
		return err
	}
	appAccessToken, err := helixClient.RequestAppAccessToken(appTokenScopes)
	if err != nil {
		return err
	}

	impl := &tokensGrpcImpl{
		globalClient: helixClient,
		appAccessToken: &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		},

		botsLock:  opts.Redsync.NewMutex("tokens-bots-lock"),
		usersLock: opts.Redsync.NewMutex("tokens-users-lock"),
		appLock:   opts.Redsync.NewMutex("tokens-app-lock"),
		db:        opts.Gorm,
		config:    opts.Config,
		log:       opts.Logger,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.TOKENS_SERVER_PORT))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	tokens.RegisterTokensServer(grpcServer, impl)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return nil
}

func (c *tokensGrpcImpl) RequestAppToken(
	_ context.Context,
	_ *emptypb.Empty,
) (*tokens.Token, error) {
	c.appLock.Lock()
	defer c.appLock.Unlock()

	if isTokenExpired(c.appAccessToken.ExpiresIn, c.appAccessToken.ObtainmentTime) {
		appAccessToken, err := c.globalClient.RequestAppAccessToken(appTokenScopes)
		if err != nil {
			return nil, err
		}

		c.appAccessToken = &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		}
		c.log.Info("app token refreshed")
	}

	return &tokens.Token{
		AccessToken: c.appAccessToken.AccessToken,
		Scopes:      []string{},
	}, nil
}

func (c *tokensGrpcImpl) RequestUserToken(
	ctx context.Context,
	data *tokens.GetUserTokenRequest,
) (*tokens.Token, error) {
	c.usersLock.Lock()
	defer c.usersLock.Unlock()

	user := model.Users{}
	err := c.db.WithContext(ctx).Where("id = ?", data.UserId).Preload("Token").Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == "" || user.Token == nil || user.Token.ID == "" {
		return nil, errors.New("cannot find user token in db")
	}

	decryptedRefreshToken, err := crypto.Decrypt(user.Token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	if isTokenExpired(int(user.Token.ExpiresIn), user.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return nil, err
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, c.config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		user.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		user.Token.AccessToken = newAccessToken

		user.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		user.Token.Scopes = newToken.Data.Scopes
		user.Token.ObtainmentTimestamp = time.Now().UTC()
		if err := c.db.WithContext(ctx).Save(&user.Token).Error; err != nil {
			return nil, err
		}
		c.log.Info("user token refreshed", slog.String("user_id", user.ID))
	}

	decryptedAccessToken, err := crypto.Decrypt(user.Token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	return &tokens.Token{
		AccessToken: decryptedAccessToken,
		Scopes:      user.Token.Scopes,
	}, nil
}

func (c *tokensGrpcImpl) RequestBotToken(
	ctx context.Context,
	data *tokens.GetBotTokenRequest,
) (*tokens.Token, error) {
	c.botsLock.Lock()
	defer c.botsLock.Unlock()

	bot := model.Bots{}
	err := c.db.WithContext(ctx).Where("id = ?", data.BotId).Preload("Token").Find(&bot).Error
	if err != nil {
		return nil, err
	}

	if bot.ID == "" || bot.Token == nil || bot.Token.ID == "" {
		return nil, errors.New("cannot find bot token in db")
	}

	decryptedRefreshToken, err := crypto.Decrypt(bot.Token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	if isTokenExpired(int(bot.Token.ExpiresIn), bot.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return nil, err
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, c.config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		bot.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		bot.Token.AccessToken = newAccessToken

		bot.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		bot.Token.Scopes = newToken.Data.Scopes
		bot.Token.ObtainmentTimestamp = time.Now().UTC()
		if err := c.db.WithContext(ctx).Save(&bot.Token).Error; err != nil {
			return nil, err
		}
		c.log.Info("bot token refreshed", slog.String("bot_id", bot.ID))
	}

	decryptedAccessToken, err := crypto.Decrypt(bot.Token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	return &tokens.Token{
		AccessToken: decryptedAccessToken,
		Scopes:      bot.Token.Scopes,
		ExpiresIn:   bot.Token.ExpiresIn,
	}, nil
}
