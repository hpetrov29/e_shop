package middleware

import (
	"encoding/json"
)

type Service interface {
	GetSession(userId string, sessionId string) (*UserClaims, error)
}

type InMemoryDb interface {
	GetKey(key string) (string, error)
}

type service struct {
	redis InMemoryDb
}

func NewMiddlewareService(a InMemoryDb) Service {
	return &service{redis: a}
}

func (s *service) GetSession(userId string, sessionId string) (*UserClaims, error) {
	keyName := "sessions:" + userId + ":" + sessionId
	result, err := s.redis.GetKey(keyName)
	claims := &UserClaims{}
	json.Unmarshal([]byte(result), claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
