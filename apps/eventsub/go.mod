module github.com/satont/twir/apps/eventsub

go 1.20

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
)

require (
	github.com/dnsge/twitch-eventsub-bindings v1.1.0
	github.com/dnsge/twitch-eventsub-framework v1.2.3
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/nicklaw5/helix/v2 v2.22.1
	github.com/redis/go-redis/v9 v9.0.5
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-20230606001334-07f3c562ff7e
	github.com/satont/twir/libs/gomodels v0.0.0-20230606001334-07f3c562ff7e
	github.com/satont/twir/libs/grpc v0.0.0-20230606001334-07f3c562ff7e
	github.com/satont/twir/libs/pubsub v0.0.0-20230606001334-07f3c562ff7e
	github.com/satont/twir/libs/twitch v0.0.0-20230606001334-07f3c562ff7e
	github.com/satont/twir/libs/types v0.0.0-20230606001334-07f3c562ff7e
	go.uber.org/zap v1.24.0
	golang.ngrok.com/ngrok v1.1.1
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.1
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/inconshreveable/log15 v3.0.0-testing.5+incompatible // indirect
	github.com/inconshreveable/log15/v3 v3.0.0-testing.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mozillazg/go-httpheader v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
)
