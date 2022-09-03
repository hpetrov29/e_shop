package users

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Service interface {
	InsertUser(user *User) (string, error)
	CreateSession(userId string, sessionId string, valueName interface{}) error
	GetPasswordFromEmail(user *UserLogin) (string, error)
	GetClaimsFromEmail(user *UserLogin) (*UserClaims, error)
	GetSession() (string, error)
}

type Rdbms interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	GetPassword(query string, values ...interface{}) (string, error)
	GetUserDetails(query string, values ...interface{}) (*UserClaims, error)
}

type InMemoryDb interface {
	GetKey(key string) (string, error)
	SetKey(key string, value interface{}, exp time.Duration) error
}

type service struct {
	mysql Rdbms
	redis InMemoryDb
}

func NewUserssService(a Rdbms, b InMemoryDb) Service {
	return &service{mysql: a, redis: b}
}

func (s *service) InsertUser(user *User) (string, error) {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	result, err := s.mysql.ExecuteQuery(query, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	lastInsertIdStr := strconv.FormatInt(lastInsertId, 10)
	return lastInsertIdStr, nil
}

func (s *service) GetPasswordFromEmail(user *UserLogin) (string, error) {
	query := "SELECT password FROM users WHERE email = ?;"
	password, err := s.mysql.GetPassword(query, user.Email)
	if err != nil {
		return "", nil
	}
	return password, nil
}

func (s *service) GetClaimsFromEmail(user *UserLogin) (*UserClaims, error) {
	query := "SELECT id AS userId, email FROM users WHERE email = ?;"
	claims, err := s.mysql.GetUserDetails(query, user.Email)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *service) CreateSession(userId string, sessionId string, valueName interface{}) error {
	key := "sessions:" + userId + ":" + sessionId
	fmt.Println(key)
	exp := time.Duration(24 * 30 * time.Hour) //TTL for the session
	err := s.redis.SetKey(key, valueName, exp)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetSession() (string, error) {
	value, err := s.redis.GetKey("key1")
	if err != nil {
		return "", err
	}
	return value, nil
}
