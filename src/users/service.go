package users

import (
	"database/sql"
	"strconv"
)

type Service interface {
	InsertUser(user *User) (string, error)
	CreateSession(userId string, sessionId string, valueName interface{}) error
}

type Rdbms interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
}

type InMemoryDb interface {
	GetKey(key string) (string, error)
	SetKey(key string, value interface{}) error
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

func (s *service) CreateSession(userId string, sessionId string, valueName interface{}) error {
	key := "sessions:" + userId + ":" + sessionId
	err := s.redis.SetKey(key, valueName)
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
