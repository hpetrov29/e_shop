package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fnmzgdt/e_shop/src/responses"
)

type Controller interface {
	Serialize(http.Handler) http.Handler
}

type middlewareController struct {
	service Service
}

func NewMIddlewareController(a InMemoryDb) Controller {
	return &middlewareController{service: &service{redis: a}}
}

func (c *middlewareController) Serialize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//refresh token - 1 year / access token - 5 mins
		refreshCookie, err := r.Cookie("refreshToken")
		if err != nil {
			fmt.Println(err.Error())
			next.ServeHTTP(w, r)
			return
		}
		accessCookie, err := r.Cookie("accessToken")
		if err != nil {
			if err.Error() == "http: named cookie not present" {
				payload, err := Validate(refreshCookie.Value)
				if err != nil {
					fmt.Println(err)
					next.ServeHTTP(w, r)
					return
				}
				sessionId := payload.(map[string]interface{})["sessionId"].(string)
				userId := payload.(map[string]interface{})["userId"].(string)
				claims, err := c.service.GetSession(userId, sessionId)
				if err != nil {
					if err.Error() == "redis: nil" {
						next.ServeHTTP(w, r)
						return
					} else {
						responses.JSONError(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				accessToken, err := NewJWT(time.Minute*5, *claims)
				if err != nil {
					responses.JSONError(w, err.Error(), http.StatusInternalServerError)
					return
				}
				accessCookie := http.Cookie{Name: "accessToken", Value: accessToken, Path: "/", Expires: time.Now().Add(time.Minute * 5), Secure: true, HttpOnly: true}
				http.SetCookie(w, &accessCookie)
				r.Header.Add("userId", userId)
				next.ServeHTTP(w, r)
				return
			} else {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		payload, _ := Validate(accessCookie.Value)
		userId := payload.(map[string]interface{})["userId"].(string)
		r.Header.Add("userId", userId)
		fmt.Println("all good")
		next.ServeHTTP(w, r)
		return
	})
}
