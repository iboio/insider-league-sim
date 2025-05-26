package builder

import (
	"database/sql"
	"fmt"

	"league-sim/config"

	_ "github.com/go-sql-driver/mysql"
)

func SqlConnectionInit() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.MySQLUser,
		config.MySQLPassword,
		config.MySQLHost,
		config.MySQLPort,
		config.MySQLDatabase)

	db, err := sql.Open("mysql", dsn)

	if err != nil {

		return nil, err
	}
	err = db.Ping()
	if err != nil {

		return nil, err
	}

	fmt.Println("MySQL connection established")

	return db, nil
}
