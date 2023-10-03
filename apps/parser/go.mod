module github.com/satont/twir/apps/parser

go 1.21

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/getsentry/sentry-go v0.23.0
	github.com/goccy/go-json v0.9.11
	github.com/google/uuid v1.3.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/imroc/req/v3 v3.41.4
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	github.com/nicklaw5/helix/v2 v2.24.0
	github.com/prometheus/client_golang v1.16.0
	github.com/redis/go-redis/v9 v9.0.5
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/gomodels v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/gopool v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/grpc v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/integrations/spotify v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/twitch v0.0.0-20230815104716-a7601217f948
	github.com/satont/twir/libs/types v0.0.0-20230815104716-a7601217f948
	github.com/satori/go.uuid v1.2.0
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/tidwall/gjson v1.16.0
	github.com/valyala/fasttemplate v1.2.2
	go.uber.org/zap v1.25.0
	golang.org/x/exp v0.0.0-20230811145659-89c5cff77bcb
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.3
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gaukas/godicttls v0.0.4 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/pprof v0.0.0-20230811205829-9131a7e9cc17 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/onsi/ginkgo/v2 v2.11.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-20 v0.3.2 // indirect
	github.com/quic-go/quic-go v0.37.4 // indirect
	github.com/refraction-networking/utls v1.4.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230815205213-6bfd019c3878 // indirect
)

replace github.com/satont/twir/libs/integrations/spotify => ../../libs/integrations/spotify

replace github.com/satont/twir/libs/config => ../../libs/config

replace github.com/satont/twir/libs/gomodels => ../../libs/gomodels

replace github.com/satont/twir/libs/types => ../../libs/types

replace github.com/satont/twir/libs/grpc => ../../libs/grpc

replace github.com/satont/twir/libs/twitch => ../../libs/twitch

replace github.com/satont/twir/libs/gopool => ../../libs/gopool
