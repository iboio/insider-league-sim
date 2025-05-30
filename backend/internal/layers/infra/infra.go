package infra

import (
	"database/sql"
	"league-sim/internal/builder"
)

type Infra struct {
	MysqlConn *sql.DB
}

func BuildInfraLayer() *Infra {
	mysqlConn, err := builder.SqlConnectionInit()
	if err != nil {
		panic(err)
	}
	return &Infra{
		MysqlConn: mysqlConn,
	}
}
