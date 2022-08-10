// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
