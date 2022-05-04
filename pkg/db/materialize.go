package db

func GenerateMzDsn(params ConnectionParams) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	return GeneratePostgresDsn(params)
}
