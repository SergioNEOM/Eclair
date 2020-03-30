package handles

import (
	"log"
	"net/http"
	//"github.com/dgrijalva/jwt-go"
)

// CheckToken middleware func for authorization
func CheckToken(nextfunc http.HandleFunc) http.HandleFunc {
	// проверяем, был ли выдан токен (и записан в куки)
	// если да, то проверить, не подошло ли время его замены
	// а потом вызываем функцию nextfunc
	log.Print("CheckToken")
	return nextfunc
}
