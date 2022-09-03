package middleware

type UserClaims struct {
	Email       string `json:"email,omitempty"`
	UserId      string `json:"userId,omitempty"`
	SessionUUID string `json:"sessionId,omitempty"`
}