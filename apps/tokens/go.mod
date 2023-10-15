module github.com/satont/twir/apps/tokens

go 1.21.0

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/crypto => ../../libs/crypto
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/logger => ../../libs/logger
)

require (
	github.com/getsentry/sentry-go v0.25.0
	github.com/go-redis/redis/v9 v9.0.0-rc.2
	github.com/go-redsync/redsync/v4 v4.7.1
	github.com/nicklaw5/helix/v2 v2.24.0
	github.com/samber/do v1.6.0
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/crypto v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/onsi/gomega v1.27.8 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230815205213-6bfd019c3878 // indirect
)
