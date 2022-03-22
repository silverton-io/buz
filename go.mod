module github.com/silverton-io/honeypot

go 1.18

require (
	cloud.google.com/go/pubsub v1.17.1
	cloud.google.com/go/storage v1.14.0
	github.com/aws/aws-sdk-go-v2 v1.14.0
	github.com/aws/aws-sdk-go-v2/config v1.13.1
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.9.1
	github.com/aws/aws-sdk-go-v2/service/firehose v1.13.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.14.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.24.1
	github.com/coocood/freecache v1.2.0
	github.com/elastic/go-elasticsearch/v8 v8.1.0
	github.com/gin-contrib/timeout v0.0.3
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/qri-io/jsonschema v0.2.1
	github.com/rs/zerolog v1.26.1
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	github.com/tidwall/gjson v1.13.0
	github.com/twmb/franz-go v1.4.0
	github.com/twmb/franz-go/pkg/kadm v0.0.0-20220301200403-ffaee5b878c6
	github.com/ulule/limiter/v3 v3.9.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.1.0 // indirect
	cloud.google.com/go/iam v0.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.14.0 // indirect
	github.com/aws/smithy-go v1.11.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twmb/franz-go/pkg/kmsg v0.0.0-20220301200403-ffaee5b878c6 // indirect
	github.com/twmb/go-rbtree v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.65.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5 // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
