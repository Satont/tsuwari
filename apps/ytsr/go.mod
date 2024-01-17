module github.com/satont/twir/apps/ytsr

go 1.21

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/sentry => ../../libs/sentry
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
)

require (
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20231218071827-5dc09a0eae99
	github.com/satont/twir/libs/logger v0.0.0-20231218071827-5dc09a0eae99
	github.com/satont/twir/libs/sentry v0.0.0-20231218071827-5dc09a0eae99
	github.com/twirapp/twir/libs/grpc v0.0.0-20231218035440-fe1a71c14ff7
	go.uber.org/fx v1.20.1
	google.golang.org/grpc v1.60.1
)

require (
	github.com/getsentry/sentry-go v0.25.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/samber/slog-common v0.13.0 // indirect
	github.com/samber/slog-multi v1.0.2 // indirect
	github.com/samber/slog-sentry/v2 v2.2.1 // indirect
	github.com/samber/slog-zerolog/v2 v2.1.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
