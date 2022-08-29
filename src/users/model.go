package users

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

type UserClaims struct {
	Email  string `json:"email,omitempty"`
	UserId string `json:"userId,omitempty"`
	SessionUUID   string `json:"sessionId,omitempty"`
}

func NewUser() User {
	now := time.Now().Unix()
	return User{CreatedAt: now}
}

func (u User) checkFields() error {
	if u.CreatedAt == 0 {
		return errors.New("CreatedAt field can't be a null value.")
	}
	if !isEmailValid(u.Email) {
		return errors.New("Please enter a valid email.")
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func (u *User) createClaims(userId string) UserClaims {
	sessionId := uuid.New().String()
	return UserClaims{Email: u.Email, UserId: userId, SessionUUID: sessionId}
}
