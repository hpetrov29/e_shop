package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fnmzgdt/e_shop/src/responses"
)

type Middleware interface {
	Serialize(http.Handler) http.Handler
	Authorize() Adapter
	GetAccessToken(w http.ResponseWriter, r *http.Request)
	AddHeader(key, value string) Adapter
	CheckMethod(method string) Adapter
	StaffAuthorize() Adapter
}

type middlewareController struct {
	service Service
}

func NewMIddlewareController(a InMemoryDb) Middleware {
	return &middlewareController{service: &service{redis: a}}
}

type Adapter func(http.Handler) http.Handler

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
		//role := payload.(map[string]interface{})["role"].(string)

		r.Header.Add("userId", userId)
		//r.Header.Add("role", role)

		next.ServeHTTP(w, r)
		return
	})
}

func (c *middlewareController) Authorize() Adapter {
	return func(next http.Handler) http.Handler {
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
}

func (c *middlewareController) AddHeader(key, value string) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
			return
		})
	}
}

func (c *middlewareController) CheckMethod(method string) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch method {
			case "POST":
				if r.Method != http.MethodPost {
					responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
					return
				}
			case "GET":
				if r.Method != http.MethodGet {
					responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
					return
				}
			case "PATCH":
				if r.Method != http.MethodPatch {
					responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
					return
				}
			case "DELETE":
				if r.Method != http.MethodDelete {
					responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
					return
				}
			default:
				responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}

func (c *middlewareController) StaffAuthorize() Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("role")
			if role != "staff" {
				responses.JSONError(w, "Action requires special staff authorization", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}
