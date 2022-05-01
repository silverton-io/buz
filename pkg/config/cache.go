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
	SchemaTable string `json:"schemaTable"`
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
