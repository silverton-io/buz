// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Http struct {
	Enabled bool `json:"enabled"`
}

type Purge struct {
	Enabled bool `json:"enabled"`
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
	RegistryTable string `json:"registryTable,omitempty"`
	// Postgres Database
	DbHost string `json:"-"`
	DbPort uint16 `json:"-"`
	DbName string `json:"-"`
	DbUser string `json:"-"`
	DbPass string `json:"-"`
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

type Registry struct {
	Backend      `json:"backend"`
	TtlSeconds   int `json:"ttlSeconds"`
	MaxSizeBytes int `json:"maxSizeBytes"`
	Purge        `json:"purge"`
	Http         `json:"http"`
}
