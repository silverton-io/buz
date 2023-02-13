// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package postgresdb

import (
	"strconv"

	"github.com/silverton-io/buz/pkg/db"
)

// GenerateDsn generates a dsn from the provided connection params
func GenerateDsn(params db.ConnectionParams) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	p := strconv.FormatUint(uint64(params.Port), 10)
	return "postgresql://" + params.User + ":" + params.Pass + "@" + params.Host + ":" + p + "/" + params.Db
}
