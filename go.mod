module github.com/danielalvarenga/snippetbox

go 1.18

// * Using `go get github.com/go-sql-driver/mysql@v1` with @v1 to add the latest v1.x.x version
// of the MySQL driver
// * Other DB drivers: https://github.com/golang/go/wiki/SQLDrivers
require github.com/go-sql-driver/mysql v1.6.0 // indirect
