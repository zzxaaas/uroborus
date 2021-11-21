package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/common/auth"
)

// Resp Resp
func Resp(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, obj)
	c.Abort()
}

// RespString RespString
func RespString(c *gin.Context, code int, s string) {
	c.JSON(code, gin.H{"message": s})
	c.Abort()
}

// checkUsername 获取用户名
func checkUsername(c *gin.Context) (username string, err error) {
	idToken, err := c.Cookie("id_token")
	if err != nil {
		RespString(c, http.StatusUnauthorized, "No token")
		return
	}
	token, claims, err := ParseToken(idToken)
	if err != nil || !token.Valid {
		RespString(c, http.StatusUnauthorized, "权限不足")
		return
	}
	username = claims.UserName
	c.Set(auth.IDTokenSubjectContextKey, username)
	return
}

func ParseToken(tokenString string) (*jwt.Token, *auth.Claims, error) {
	Claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return auth.Jwtkey, nil
	})
	return token, Claims, err
}

// Auth middleware
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := checkUsername(c)
		if username == "" {
			return
		}

		c.Next()
	}
}
