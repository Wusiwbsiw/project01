module https_operation

go 1.25.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	sql_operation v0.0.0-unspecified
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
)

replace sql_operation => ../sql_operation
