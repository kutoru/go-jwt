package tokens

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kutoru/go-jwt/glb"
)

// Creates a cookie with a new access token as a value
func CreateAccessCookie(guid string, exp time.Time) (*http.Cookie, error) {
	claims := jwt.MapClaims{}
	claims["guid"] = guid
	claims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "access",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   glb.EXP,
		SameSite: http.SameSiteNoneMode,
	}

	// log.Println("Created access token:", tokenString)

	return cookie, nil
}

// Parses the access token from the requests' cookies and returns its GUID and EXP unix time if it is valid. Otherwise returns an error
func GetAccessTokenInfo(r *http.Request) (string, int64, error) {
	cookie, err := r.Cookie("access")
	if err != nil {
		return "", 0, err
	}

	token, err := jwt.Parse(cookie.Value, parseHelper)
	if err != nil {
		return "", 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, fmt.Errorf("token claims are invalid")
	}

	guid, ok := claims["guid"].(string)
	if !ok {
		return "", 0, fmt.Errorf("guid is invalid")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("exp is invalid")
	}

	if int64(exp) <= time.Now().Unix() {
		return "", 0, fmt.Errorf("token has expired")
	}

	return guid, int64(exp), nil
}

func parseHelper(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
	}

	return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
}
