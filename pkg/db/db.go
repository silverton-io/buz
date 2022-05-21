package db

const (
	POSTGRES      = "postgres"
	MATERIALIZE   = "materialize"
	MYSQL         = "mysql"
	CLICKHOUSE    = "clickhouse"
	MONGODB       = "mongodb"
	ELASTICSEARCH = "elasticsearch"
	TIMESCALE     = "timescale"
)

type ConnectionParams struct {
	Host string
	Port uint16
	Db   string
	User string
	Pass string
}
