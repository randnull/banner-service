package jwt_token

import (
	"fmt"
	"log"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/banner-service/internal/errors"
)


func CreateJWTToken(is_admin bool, jwt_secret string) (string, error) {
	secret_key := []byte(jwt_secret)

	claims := jwt.MapClaims{
		"is_admin": is_admin,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	token_str, err := token.SignedString(secret_key)

	if err != nil {
		log.Fatal(err)
	}

	return token_str, nil
}


func ParseJWTToken(token_str string, jwt_secret string) (bool, error) {
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_secret), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.InvalidToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	str := fmt.Sprintf("%v", claims["is_admin"])

	id, _ := strconv.ParseBool(str)
	
	return id, nil
}
