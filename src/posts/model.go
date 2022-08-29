package posts

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserId    string `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

func NewPost(userId string) Post {
	now := time.Now().Unix()
	return Post{UserId: userId, CreatedAt: now}
}

func (p Post) checkFields() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("Title field can't be empty.")
	}
	if strings.TrimSpace(p.Body) == "" {
		return errors.New("Body field can't be empty.")
	}
	if p.CreatedAt == 0 {
		return errors.New("CreatedAt field can't be a null value.")
	}
	if p.UserId == "" {
		return errors.New("UserId field can't be a null value.")
	}
	return nil
}
