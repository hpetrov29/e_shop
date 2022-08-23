package posts

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserId    int    `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

func NewPost(userId int) Post {
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
	if p.UserId == 0 {
		return errors.New("UserId field can't be a null value.")
	}
	return nil
}
