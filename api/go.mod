module github.com/and-period/furumaru/api

go 1.22

require (
	firebase.google.com/go/v4 v4.14.1
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/config v1.27.26
	github.com/aws/aws-sdk-go-v2/credentials v1.17.26
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.14.9
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.41.4
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.3
	github.com/aws/aws-sdk-go-v2/service/mediaconvert v1.57.3
	github.com/aws/aws-sdk-go-v2/service/medialive v1.54.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.58.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.3
	github.com/aws/aws-sdk-go-v2/service/sfn v1.29.3
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.3
	github.com/casbin/casbin/v2 v2.97.0
	github.com/dlclark/regexp2 v1.11.2
	github.com/getsentry/sentry-go v0.28.1
	github.com/gin-contrib/gzip v0.0.6
	github.com/gin-contrib/pprof v1.5.0
	github.com/gin-contrib/zap v1.1.3
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/jinzhu/copier v0.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/line/line-bot-sdk-go/v7 v7.21.0
	github.com/newrelic/go-agent/v3 v3.33.1
	github.com/newrelic/go-agent/v3/integrations/nrgin v1.3.1
	github.com/prometheus/client_golang v1.19.1
	github.com/rafaelhl/gorm-newrelic-telemetry-plugin v1.0.0
	github.com/rs/cors v1.11.0
	github.com/rs/cors/wrapper/gin v0.0.0-20240228164225-8d33ca4794ea
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.14.0+incompatible
	github.com/shopspring/decimal v1.4.0
	github.com/slack-go/slack v0.13.1
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.9.0
	github.com/stripe/stripe-go/v73 v73.16.0
	go.uber.org/mock v0.4.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.7.0
	google.golang.org/api v0.188.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240708141625-4ad9e859172b
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.2.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.11
	moul.io/zapgorm2 v1.3.0
)

require (
	cloud.google.com/go/auth v0.7.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.3 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
)

require (
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	cloud.google.com/go/firestore v1.15.0 // indirect
	cloud.google.com/go/iam v1.1.11 // indirect
	cloud.google.com/go/longrunning v0.5.10 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/ivs v1.37.3
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.3 // indirect
	github.com/aws/smithy-go v1.20.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.11.9 // indirect
	github.com/casbin/govaluate v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/oauth2 v0.21.0
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/appengine/v2 v2.0.6 // indirect
	google.golang.org/genproto v0.0.0-20240708141625-4ad9e859172b // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240701130421-f6361c86f094 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
