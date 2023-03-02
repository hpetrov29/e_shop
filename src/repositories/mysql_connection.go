package repositories

import (
	"database/sql"
	"fmt"

	"github.com/fnmzgdt/e_shop/src/items"
	"github.com/fnmzgdt/e_shop/src/users"
	"github.com/fnmzgdt/e_shop/src/utils"
	_ "github.com/go-sql-driver/mysql"
)

type SqlConnection struct {
	db *sql.DB
}

func SetupSqlConnection() (*SqlConnection, error) {
	var (
		dbname   = utils.GetEnv("MYSQL_DB_NAME", "")
		user     = utils.GetEnv("MYSQL_USER", "root")
		password = utils.GetEnv("MYSQL_PASSWORD", "")
		host     = utils.GetEnv("MYSQL_HOST", "localhost")
	)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)) //host = host.docker.internal for docker dev
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successful conneciton to MySQL.")
	return &SqlConnection{db: db}, nil
}

func (s *SqlConnection) ExecuteQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SqlConnection) GetPassword(query string, values ...interface{}) (string, error) {
	var password string

	err := s.db.QueryRow(query, values...).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *SqlConnection) GetUserDetails(query string, values ...interface{}) (*users.UserClaims, error) {
	userClaims := users.UserClaims{}
	err := s.db.QueryRow(query, values...).Scan(&userClaims.UserId, &userClaims.Email)
	if err != nil {
		return nil, err
	}
	return &userClaims, nil
}

func (s *SqlConnection) GetItem(query string, id int) (*items.ItemGet, error) {
	item := items.ItemGet{}
	if err := s.db.QueryRow(query, id).Scan(&item.Id, &item.UserId, &item.CategoryName, &item.BrandName, &item.CreatedAt, &item.Price, &item.DiscountedPrice, &item.Description, &item.ModifiedAt); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SqlConnection) GetImages(query string, itemId int) (*[]items.Image, error) {
	imagesArray := make([]items.Image, 0)
	rows, err := s.db.Query(query, itemId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		image := new(items.Image)
		if err := rows.Scan(&image.Url); err != nil {
			return nil, err
		}
		imagesArray = append(imagesArray, *image)
	}
	return &imagesArray, nil
}

func (s *SqlConnection) GetItems(query string, values ...interface{}) (*[]items.ItemGet, error) {
	itemsArray := make([]items.ItemGet, 0)
	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := new(items.ItemGet)
		item.Images = make([]items.Image, 1)
		if err := rows.Scan(&item.Id, &item.UserId, &item.CategoryName, &item.BrandName, &item.CreatedAt, &item.Price, &item.DiscountedPrice, &item.Description, &item.ModifiedAt, &item.Images[0].Url); err != nil {
			return nil, err
		}
		itemsArray = append(itemsArray, *item)
	}
	return &itemsArray, nil
}
