package utils

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
    hmacSecretString := os.Getenv("SECRET_JWT")
    hmacSecret := []byte(hmacSecretString)
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return hmacSecret, nil
    })

    if err != nil {
		log.Println(err.Error())

        return nil, false
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, true
    } else {
        log.Printf("Invalid JWT Token")
        return nil, false
    }
}