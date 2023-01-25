module github.com/and-period/furumaru/api

go 1.19

require (
	firebase.google.com/go/v4 v4.10.0
	github.com/aws/aws-lambda-go v1.36.0
	github.com/aws/aws-sdk-go-v2 v1.17.3
	github.com/aws/aws-sdk-go-v2/config v1.18.7
	github.com/aws/aws-sdk-go-v2/credentials v1.13.7
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.8
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.21.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.17.9
	github.com/aws/aws-sdk-go-v2/service/s3 v1.29.6
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.16.5
	github.com/aws/aws-sdk-go-v2/service/sqs v1.20.0
	github.com/casbin/casbin/v2 v2.60.0
	github.com/gin-contrib/gzip v0.0.6
	github.com/gin-contrib/zap v0.1.0
	github.com/gin-gonic/gin v1.8.2
	github.com/go-playground/validator/v10 v10.11.1
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.3
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/line/line-bot-sdk-go/v7 v7.18.0
	github.com/newrelic/go-agent/v3 v3.18.2
	github.com/newrelic/go-agent/v3/integrations/nrgin v1.1.3
	github.com/prometheus/client_golang v1.14.0
	github.com/rafaelhl/gorm-newrelic-telemetry-plugin v1.0.0
	github.com/rs/cors v1.8.3
	github.com/rs/cors/wrapper/gin v0.0.0-20220619195839-da52b0701de5
	github.com/samber/lo v1.37.0
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.12.0+incompatible
	github.com/shopspring/decimal v1.3.1
	github.com/slack-go/slack v0.11.4
	github.com/stretchr/testify v1.8.1
	github.com/stripe/stripe-go/v73 v73.4.0
	go.uber.org/zap v1.24.0
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/sync v0.1.0
	google.golang.org/api v0.104.0
	google.golang.org/genproto v0.0.0-20221206210731-b1a01be3a5f6
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.1.0
	gorm.io/driver/mysql v1.4.4
	gorm.io/gorm v1.24.2
	moul.io/zapgorm2 v1.2.0
)

require (
	github.com/MicahParks/keyfunc v1.5.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
)

require (
	cloud.google.com/go v0.105.0 // indirect
	cloud.google.com/go/compute v1.13.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.2 // indirect
	cloud.google.com/go/firestore v1.9.0 // indirect
	cloud.google.com/go/iam v0.8.0 // indirect
	cloud.google.com/go/longrunning v0.3.0 // indirect
	cloud.google.com/go/storage v1.27.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.13.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/ivs v1.19.0
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.0 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20221005025214-4161e89ecf1b // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/oauth2 v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/appengine/v2 v2.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
