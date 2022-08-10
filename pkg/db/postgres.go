// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package db

import "strconv"

func GeneratePostgresDsn(params ConnectionParams) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	p := strconv.FormatUint(uint64(params.Port), 10)
	return "postgresql://" + params.User + ":" + params.Pass + "@" + params.Host + ":" + p + "/" + params.Db
}
