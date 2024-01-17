package workflow

import (
	"time"

	"github.com/satont/twir/apps/timers/internal/activity"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	"github.com/satont/twir/apps/timers/internal/shared"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Lc     fx.Lifecycle
	Logger logger.Logger
	Cfg    cfg.Config

	TimersRepository   timers.Repository
	ChannelsRepository channels.Repository
	StreamsRepository  streams.Repository

	BotsGrpc   bots.BotsClient
	ParserGrpc parser.ParserClient
	Activity   *activity.Activity
}

func New(opts Opts) (*Workflow, error) {
	cl, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Workflow{
		logger:             opts.Logger,
		config:             opts.Cfg,
		cl:                 cl,
		timersRepository:   opts.TimersRepository,
		channelsRepository: opts.ChannelsRepository,
		streamsRepository:  opts.StreamsRepository,
		botsGrpc:           opts.BotsGrpc,
		parserGrpc:         opts.ParserGrpc,
		activity:           opts.Activity,
	}, nil
}

type Workflow struct {
	logger logger.Logger
	config cfg.Config
	db     *gorm.DB
	cl     client.Client

	timersRepository   timers.Repository
	channelsRepository channels.Repository
	streamsRepository  streams.Repository

	botsGrpc   bots.BotsClient
	parserGrpc parser.ParserClient
	activity   *activity.Activity
}

func (c *Workflow) Flow(ctx workflow.Context, timer timers.Timer) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		HeartbeatTimeout:    time.Second * 10,
		TaskQueue:           shared.TimersWorkerTaskQueueName,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Second * 100,
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	workflowInfo := workflow.GetInfo(ctx)
	lastResponseValue := workflowInfo.Memo.Fields["lastResponse"]
	var currentResponse int
	err := converter.GetDefaultDataConverter().FromPayload(lastResponseValue, &currentResponse)
	if err != nil {
		return err
	}

	var newResponse int
	err = workflow.ExecuteActivity(
		ctx,
		c.activity.SendMessage,
		timer.ID,
		currentResponse,
	).Get(ctx, &newResponse)
	if err != nil {
		return err
	}

	memo := map[string]interface{}{
		"lastResponse": newResponse,
	}

	err = workflow.UpsertMemo(ctx, memo)
	if err != nil {
		return err
	}

	return nil
}
