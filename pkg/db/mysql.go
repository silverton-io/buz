// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package db

import "strconv"

func GenerateMysqlDsn(params ConnectionParams) string {
	port := strconv.FormatUint(uint64(params.Port), 10)
	return params.User + ":" + params.Pass + "@tcp(" + params.Host + ":" + port + ")/" + params.Db + "?parseTime=true"
}
