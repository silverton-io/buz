// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Purge struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type Backend struct {
	Type string `json:"type"`
	Path string `json:"path"`
	// S3 and Gcs
	Bucket string `json:"bucket,omitempty"`
	// Gcs
	Region string `json:"region,omitempty"`
	// Http
	Host string `json:"host,omitempty"`
	// Db, general
	RegistryTable string `json:"registryTable"`
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
	// Clickhouse Database
	ClickhouseHost   string `json:"-"`
	ClickhousePort   uint16 `json:"-"`
	ClickhouseDbName string `json:"-"`
	ClickhouseUser   string `json:"-"`
	ClickhousePass   string `json:"-"`
	// Mongodb
	MongoHosts         []string `json:"mongoHosts,omitempty"`
	MongoPort          string   `json:"mongoDbPort,omitempty"`
	MongoDbName        string   `json:"mongoDbName,omitempty"`
	MongoUser          string   `json:"-"`
	MongoPass          string   `json:"-"`
	RegistryCollection string   `json:"registryCollection,omitempty"`
	// Minio
	MinioEndpoint   string `json:"minioEndpoint,omitempty"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

type SchemaDirectory struct {
	Enabled bool `json:"enabled"`
}

type SchemaCache struct {
	Backend         `json:"backend"`
	TtlSeconds      int `json:"ttlSeconds"`
	MaxSizeBytes    int `json:"maxSizeBytes"`
	Purge           `json:"purge"`
	SchemaDirectory `json:"schemaDirectory"`
}
