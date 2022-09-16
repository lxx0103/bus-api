package service

import (
	"errors"

	"bus-api/core/config"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserID   int64
	UserName string
	Role     string
	jwt.StandardClaims
}

//jwt service
type JWTService interface {
	GenerateToken(claims CustomClaims) string
	ParseToken(tokenString string) (*CustomClaims, error)
}

type jwtServices struct {
	secretKey []byte
}

//auth-jwt
func JWTAuthService() JWTService {
	jwtSecret := config.ReadConfig("auth.secret")
	return &jwtServices{
		secretKey: []byte(jwtSecret),
	}
}

func (service *jwtServices) GenerateToken(claims CustomClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// 解析 token
func (service *jwtServices) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return service.secretKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid1")
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("TokenInvalid2")

	} else {
		return nil, errors.New("TokenInvalid3")

	}

}
