package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/satont/tsuwari/apps/parser/internal/commands"
	"github.com/satont/tsuwari/apps/parser/internal/config/redis"
	"github.com/satont/tsuwari/apps/parser/internal/variables"

	cfg "github.com/satont/tsuwari/libs/config"

	twitch "github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	usersauth "github.com/satont/tsuwari/apps/parser/internal/twitch/user"

	"github.com/getsentry/sentry-go"

	"github.com/satont/tsuwari/apps/parser/internal/grpc_impl"
	"github.com/satont/tsuwari/libs/grpc/clients"
	parser "github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	if cfg.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	r := redis.New(cfg.RedisUrl)
	defer r.Close()

	botsGrpcClient := clients.NewBots(cfg.AppEnv)
	dotaGrpcClient := clients.NewDota(cfg.AppEnv)
	evalGrpcClient := clients.NewEval(cfg.AppEnv)

	usersAuthService := usersauth.New(usersauth.UsersServiceOpts{
		Db:           db,
		ClientId:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	})
	twitchClient := twitch.New(*cfg)
	variablesService := variables.New()
	commandsService := commands.New(commands.CommandsOpts{
		Redis:            r,
		VariablesService: variablesService,
		Db:               db,
		UsersAuth:        usersAuthService,
		Twitch:           twitchClient,
		BotsGrpc:         botsGrpcClient,
		DotaGrpc:         dotaGrpcClient,
		EvalGrpc:         evalGrpcClient,
	})

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.PARSER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	parser.RegisterParserServer(grpcServer, grpc_impl.NewServer(&grpc_impl.GrpcImplOpts{
		Redis:     r,
		Variables: &variablesService,
		Commands:  &commandsService,
	}))
	go grpcServer.Serve(lis)

	logger.Info("Started")

	// runtime.Goexit()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
	grpcServer.Stop()
	d.Close()
}
