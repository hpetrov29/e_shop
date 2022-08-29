package posts

import "database/sql"

type Service interface {
	WritePost(post *Post) error
}

type Rdbms interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
}

type service struct {
	mysql Rdbms
}

func NewPostsService(db Rdbms) Service {
	return &service{db}
}

func (s *service) WritePost(post *Post) error {
	query := "INSERT INTO posts(title, body) VALUES (?, ?)"
	_, err := s.mysql.ExecuteQuery(query, post.Title, post.Body)
	if err != nil {
		return err
	}
	return nil
}
