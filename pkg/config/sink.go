package config

type Sink struct {
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	DeliveryRequired bool     `json:"deliveryRequired"`
	Project          string   `json:"project,omitempty"`
	KafkaBrokers     []string `json:"kakfaBrokers,omitempty"`
	// Kafka, Pubsub
	ValidTopic   string `json:"validTopic,omitempty"`
	InvalidTopic string `json:"invalidTopic,omitempty"`
	// Kinesis
	ValidStream   string `json:"validStream,omitempty"`
	InvalidStream string `json:"invalidStream,omitempty"`
	// Relay, HTTP/S, etc
	ValidUrl   string `json:"validUrl,omitempty"`
	InvalidUrl string `json:"invalidUrl,omitempty"`
	// Subject-based
	ValidSubject   string `json:"validSubject,omitempty"`
	InvalidSubject string `json:"invalidSubject,omitempty"`
	// NATS
	NatsHost string `json:"-"`
	NatsUser string `json:"-"`
	NatsPass string `json:"-"`
	// Elasticsearch
	ValidIndex            string   `json:"validIndex,omitempty"`
	InvalidIndex          string   `json:"invalidIndex,omitempty"`
	ElasticsearchHosts    []string `json:"elasticsearchHosts,omitempty"`
	ElasticsearchUsername string   `json:"-"`
	ElasticsearchPassword string   `json:"-"`
	// Honeypot relay
	RelayUrl string `json:"relayUrl,omitempty"`
	// File
	ValidFile   string `json:"validFile,omitempty"`
	InvalidFile string `json:"invalidFile,omitempty"`
	// Postgres Database
	PgHost   string `json:"-"`
	PgPort   uint16 `json:"-"`
	PgDbName string `json:"-"`
	PgUser   string `json:"-"`
	PgPass   string `json:"-"`
	// Mysql Database
	MysqlHost   string `json:"-"`
	MysqlPort   uint16 `json:"-"`
	MysqlDbName string `json:"-"`
	MysqlUser   string `json:"-"`
	MysqlPass   string `json:"-"`
	// Materialize Database
	MzHost   string `json:"-"`
	MzPort   uint16 `json:"-"`
	MzDbName string `json:"-"`
	MzUser   string `json:"-"`
	MzPass   string `json:"-"`
	// Timescale Database
	TimescaleHost   string `json:"-"`
	TimescalePort   uint16 `json:"-"`
	TimescaleDbName string `json:"-"`
	TimescaleUser   string `json:"-"`
	TimescalePass   string `json:"-"`
	// Clickhouse Database
	ClickhouseHost   string `json:"-"`
	ClickhousePort   uint16 `json:"-"`
	ClickhouseDbName string `json:"-"`
	ClickhouseUser   string `json:"-"`
	ClickhousePass   string `json:"-"`
	// Database
	ValidTable   string `json:"validTable,omitempty"`
	InvalidTable string `json:"invalidTable,omitempty"`
	// Pubnub
	ValidChannel   string `json:"validChannel,omitempty"`
	InvalidChannel string `json:"invalidChannel,omitempty"`
	PubnubPubKey   string `json:"pubnubPubKey,omitempty"`
	PubnubSubKey   string `json:"pubnubSubKey,omitempty"`
	// Mongodb
	MongoHosts        []string `json:"mongoHosts,omitempty"`
	MongoPort         string   `json:"mongoDbPort,omitempty"`
	MongoDbName       string   `json:"mongoDbName,omitempty"`
	MongoUser         string   `json:"-"`
	MongoPass         string   `json:"-"`
	ValidCollection   string   `json:"validCollection,omitempty"`
	InvalidCollection string   `json:"invalidCollection,omitempty"`
	// Indicative
	IndicativeApiKey string `json:"-"`
	// Amplitude
	AmplitudeApiKey string `json:"-"`
}
