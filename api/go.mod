module github.com/and-period/furumaru/api

go 1.18

require (
	github.com/aws/aws-lambda-go v1.32.0
	github.com/aws/aws-sdk-go-v2 v1.16.6
	github.com/aws/aws-sdk-go-v2/config v1.15.9
	github.com/aws/aws-sdk-go-v2/credentials v1.12.4
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.16.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.10
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.15.12
	github.com/aws/aws-sdk-go-v2/service/sqs v1.18.6
	github.com/casbin/casbin/v2 v2.47.1
	github.com/gin-contrib/gzip v0.0.5
	github.com/gin-contrib/zap v0.0.2
	github.com/gin-gonic/gin v1.8.0
	github.com/go-playground/validator/v10 v10.11.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.2
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/line/line-bot-sdk-go/v7 v7.16.0
	github.com/newrelic/go-agent/v3 v3.17.0
	github.com/newrelic/go-agent/v3/integrations/nrgin v1.1.2
	github.com/prometheus/client_golang v1.12.2
	github.com/rafaelhl/gorm-newrelic-telemetry-plugin v1.0.0
	github.com/rs/cors v1.8.2
	github.com/rs/cors/wrapper/gin v0.0.0-20220223021805-a4a5ce87d5a2
	github.com/samber/lo v1.26.0
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.11.1+incompatible
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.7.1
	go.uber.org/zap v1.21.0
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29
	google.golang.org/genproto v0.0.0-20220531134929-86cf59382f1b
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.6
	gorm.io/driver/mysql v1.3.4
	gorm.io/gorm v1.23.5
	moul.io/zapgorm2 v1.1.3
)

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.6 // indirect
	github.com/aws/smithy-go v1.12.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
