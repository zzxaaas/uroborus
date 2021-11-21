package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	IDTokenSubjectContextKey = "id_token_sub"
	IDTokenSubjectKey        = "id_token"
)

var Jwtkey = []byte("www.topgoer.com")

type Claims struct {
	UserName string
	jwt.StandardClaims
}

//颁发token
func SetToken(userName string) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(Jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
