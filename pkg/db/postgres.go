package db

import "strconv"

func GeneratePostgresDsn(params DbConnectionParams) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	p := strconv.FormatUint(uint64(params.Port), 10)
	return "postgresql://" + params.User + ":" + params.Pass + "@" + params.Host + ":" + p + "/" + params.Db
}
