package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fnmzgdt/e_shop/src/responses"
)

type Controller interface {
	Serialize(http.Handler) http.Handler
	Authorize(next http.Handler) http.Handler
	GetAccessToken(w http.ResponseWriter, r *http.Request)
}

type middlewareController struct {
	service Service
}

func NewMIddlewareController(a InMemoryDb) Controller {
	return &middlewareController{service: &service{redis: a}}
}

func (c *middlewareController) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refreshToken")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "http://google.com", 302)
		return
	}
	refreshJwt := refreshCookie.Value
	payload, err := Validate(refreshJwt)
	if err != nil {
		if err.Error() == "validate: Token is expired" {
			http.Redirect(w, r, "http://google.com", 302)
			return
		}
		responses.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionId := payload.(map[string]interface{})["sessionId"].(string)
	userId := payload.(map[string]interface{})["userId"].(string)
	claims, err := c.service.GetSession(userId, sessionId)
	if err != nil {
		if err.Error() == "redis: nil" {
			http.Redirect(w, r, "http://google.com", 302)
			//delete cookie
			return
		} else {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//check if session.userId == userId
	accessToken, err := NewJWT(time.Minute*5, *claims)
	if err != nil {
		responses.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accessCookie := http.Cookie{Name: "accessToken", Value: accessToken, Path: "/", Expires: time.Now().Add(time.Minute * 5), Secure: true, HttpOnly: true}
	http.SetCookie(w, &accessCookie)
	responses.JSONResponse(w, "Access token successfully renewed.", nil, 200)
	return
}

func (c *middlewareController) Serialize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessCookie, err := r.Cookie("accessToken")
		if err != nil {
			if err.Error() == "http: named cookie not present" {
				next.ServeHTTP(w, r)
				return
			} else {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		payload, err := Validate(accessCookie.Value)
		if err != nil {
			if err.Error() == "validate: Token is expired" {
				next.ServeHTTP(w, r)
				return
			}
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId := payload.(map[string]interface{})["userId"].(string)
		r.Header.Add("userId", userId)
		next.ServeHTTP(w, r)
		return
	})
}

func (c *middlewareController) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		if userId == "" {
			responses.JSONError(w, "Action requires authorization", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
