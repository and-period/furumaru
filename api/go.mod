module github.com/and-period/furumaru/api

go 1.24

require (
	firebase.google.com/go/v4 v4.15.2
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.9
	github.com/aws/aws-sdk-go-v2/credentials v1.17.62
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.18.7
	github.com/aws/aws-sdk-go-v2/service/batch v1.51.1
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.51.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.41.1
	github.com/aws/aws-sdk-go-v2/service/mediaconvert v1.69.1
	github.com/aws/aws-sdk-go-v2/service/medialive v1.71.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.78.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.2
	github.com/aws/aws-sdk-go-v2/service/sfn v1.35.1
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.1
	github.com/aws/aws-sdk-go-v2/service/transcribe v1.44.2
	github.com/aws/aws-sdk-go-v2/service/translate v1.29.1
	github.com/casbin/casbin/v2 v2.103.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/dlclark/regexp2 v1.11.5
	github.com/getsentry/sentry-go v0.31.1
	github.com/gin-contrib/gzip v1.2.2
	github.com/gin-contrib/pprof v1.5.2
	github.com/gin-contrib/zap v1.1.4
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.25.0
	github.com/go-sql-driver/mysql v1.9.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.1
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/jinzhu/copier v0.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/line/line-bot-sdk-go/v7 v7.21.0
	github.com/newrelic/go-agent/v3 v3.37.0
	github.com/newrelic/go-agent/v3/integrations/nrgin v1.3.3
	github.com/prometheus/client_golang v1.21.1
	github.com/rafaelhl/gorm-newrelic-telemetry-plugin v1.0.0
	github.com/rs/cors v1.11.1
	github.com/rs/cors/wrapper/gin v0.0.0-20240830163046-1084d89a1692
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.16.0+incompatible
	github.com/shopspring/decimal v1.4.0
	github.com/slack-go/slack v0.16.0
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	github.com/stripe/stripe-go/v73 v73.16.0
	go.uber.org/goleak v1.3.0
	go.uber.org/mock v0.5.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.12.0
	google.golang.org/api v0.225.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	googlemaps.github.io/maps v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.2.5
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
	moul.io/zapgorm2 v1.3.0
)

require (
	cel.dev/expr v0.22.0 // indirect
	cloud.google.com/go v0.119.0 // indirect
	cloud.google.com/go/auth v0.15.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/firestore v1.18.0 // indirect
	cloud.google.com/go/iam v1.4.1 // indirect
	cloud.google.com/go/longrunning v0.6.5 // indirect
	cloud.google.com/go/monitoring v1.24.0 // indirect
	cloud.google.com/go/storage v1.50.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.27.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.51.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.51.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/ivs v1.43.1
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/bytedance/sonic v1.13.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cncf/xds/go v0.0.0-20250121191232-2f005788dc42 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.32.4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.5 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.35.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.60.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/oauth2 v0.28.0
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/appengine/v2 v2.0.6 // indirect
	google.golang.org/genproto v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
