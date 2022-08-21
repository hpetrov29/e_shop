package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnection struct {
	db *sql.DB
}

func SetupMySQLConnection() (*MySQLConnection, error) {
	db, err := sql.Open("mysql", "root:+Zrtp2B&Eur27@/go_chat_app")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successful conneciton to MySQL.")
	return &MySQLConnection{db: db}, nil
}

func (s *MySQLConnection) ExecuteQuery(query string, values ...interface{}) (error) {
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}