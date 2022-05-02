package db

import "strconv"

func GenerateMysqlDsn(params ConnectionParams) string {
	port := strconv.FormatUint(uint64(params.Port), 10)
	return params.User + ":" + params.Pass + "@tcp(" + params.Host + ":" + port + ")/" + params.Db + "?parseTime=true"
}
