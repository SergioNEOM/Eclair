package handles

import (
	"errors"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

const SECRET = "Eclair-#&#-Secret"

// JWTClaims is a struct Public part of claims.
type JWTCustomClaims struct {
	Login string `json:"login"`
	Role  int    `json:"role"`
	jwt.StandardClaims
}

// CheckToken verify token string
func CheckToken(tokenStr string) error {

	// получить строку токена и расшифровать ее
	// если да, то проверить, не подошло ли время его замены
	// если время замены не просрочено, выдать новый токен и записать в куки
	// иначе - отправить на авторизацию
	log.Print("CheckToken begin")
	//fmt.Println("CheckToken ---")
	if tokenStr == "" {
		// не был выдан ??
		// уходим на выдачу...
		log.Print("token not found - create new")
	}

	t, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
	if err != nil {
		log.Println("Token parse error. Go to authentification")
		//todo: redirect ????
		return err
	}

	if claims, ok := t.Claims.(*JWTCustomClaims); ok && t.Valid {
		log.Printf("%v %v %v\n", claims.Login, claims.Role, claims.StandardClaims.ExpiresAt)
	} else {
		log.Println(err)
		return errors.New("Token isn't recognized")
	}

	return nil
}

func CheckTokenAuth() {
	log.Print("CheckTokenAuth")
	fmt.Println("CheckToken --- auth")
}
