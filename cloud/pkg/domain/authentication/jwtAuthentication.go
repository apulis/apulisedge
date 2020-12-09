package authentication

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func jwtAuthtication(c *gin.Context) {
	r := c.Request
	auth := r.Header.Get("Authorization")
	fmt.Println(auth)
	tokenString := strings.Fields(auth)[1]
	fmt.Println("token acquired: ", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("Sign key for JWT"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		c.Abort()
		fmt.Println(err)
	}
	c.Next()
}

func jwtSign() string {
	mySigningKey := []byte("Sign key for JWT")

	type MyCustomClaims struct {
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Date(2200, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	return "success"
}
