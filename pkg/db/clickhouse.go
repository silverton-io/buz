// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package db

import "strconv"

// GenerateClickhouseDsn generates a dsn from the provided connection params
func GenerateClickhouseDsn(params ConnectionParams) string {
	// "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
	port := strconv.FormatUint(uint64(params.Port), 10)
	return "tcp://" + params.Host + ":" + port + "?database=" + params.Db + "&username=" + params.User + "&password=" + params.Pass
}
