package handles

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const SECRET = "Eclair-#&#-Secret"

// JWTCustomClaims is a struct Public part of claims.
type JWTCustomClaims struct {
	Login string `json:"login"`
	Role  int    `json:"role"`
	jwt.StandardClaims
}

//
//CheckJWToken проверяет валидность пользователя
func CheckJWToken(w http.ResponseWriter, r *http.Request) {
	// 1. получить куку. Если нет (или пустое значение) -
	// сверяем логин в БД. Если там нет, то редирект на  "/auth", если есть, выдаем новый токен
	tokenStr := getCookieValue(r)
	if tokenStr == "" {
		log.Printf("Request %v - Token not found in cookie. Go to authentification", r.RequestURI)
		/////
		http.Redirect(w, r, "/", 307 /*401*/)
		return
	}
	// 2. Получить в БД refreshToken, чтобы на нём проверить accessToken.
	//  если не парсится, считаем скопрометированным - редирект на  "/auth"
	secretKey := SECRET //todo: потом заменить функцией получения из БД
	//
	t, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Printf("Token parse error. Go to authentification --- %v\n", err)
		http.Redirect(w, r, "/", 401)
		return
	}

	// 3. после парсинга проверить время жизни. если истекло - редирект на refreshJWTPath
	claims, ok := t.Claims.(*JWTCustomClaims)
	if !(ok && t.Valid) {
		log.Println("Token parse error - claims not recognized. Go to authentification")
		http.Redirect(w, r, "/", 401)
		return
	}

	if time.Now().After(time.Unix(claims.StandardClaims.ExpiresAt, 0)) {
		log.Printf("Token expired date: %v %v %v\n", claims.Login, claims.Role, claims.StandardClaims.ExpiresAt)
		http.Redirect(w, r, refreshJWTPath, 417)
		return
	}
	// 4. если всё норм, можем проверить поля токена (можем сверить роль с БД) ???????

	// 5. возврат в основной поток
	return
}

//==================

// CreateJWT create new JWT
func CreateJWT(claims *JWTCustomClaims, expirationTime *time.Time, jwtKey []byte) (string, error) {
	// JWTExpiresAt больше 0 установим expiry time
	if expirationTime != nil {
		claims.StandardClaims.ExpiresAt = expirationTime.Unix()
	}

	//log.Println("JWT: ", claims)

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	//log.Println("JWT: ", tokenString)

	return tokenString, nil
}

/*
// CreateJWTCookie create new JWT as Cookie
func CreateJWTCookie(claims *JWTCustomClaims, jwtExpiresAt int, jwtKey []byte) (*http.Cookie, error) {
	var expirationTime *time.Time

	// jwtExpiresAt > 0 установим expiry time
	if jwtExpiresAt > 0 {
		t := time.Now().Add(time.Duration(jwtExpiresAt * int(time.Second)))
		expirationTime = &t
	}

	// создадим новый токен
	tokenString, err := CreateJWT(claims, expirationTime, jwtKey)
	if err != nil {
		return nil, err
	}

	// подготовим Cookie
	cookie := http.Cookie{
		Name:  SECRET_COOKIE,
		Value: tokenString,
	}

	if jwtExpiresAt > 0 {
		// set an expiry time is the same as the token itself
		cookie.Expires = *expirationTime
	} else {
		cookie.MaxAge = 0 // без ограничения времени жизни
	}

	return &cookie, nil
}
*/
