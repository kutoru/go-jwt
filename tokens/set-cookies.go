package tokens

import (
	"log"
	"net/http"
	"time"

	"github.com/kutoru/go-jwt/glb"
)

// Tries to create both access and refresh tokens and set them to the writer as cookies
func CreateAndSetAsCookies(w http.ResponseWriter, guid int) error {
	// Getting expiry time
	exp := time.Now().Add(time.Duration(glb.EXP) * time.Second)

	// Creating a cookie that contains the access token
	accessCookie, err := CreateAccessCookie(guid, exp)
	if err != nil {
		return err
	}

	// Creating a cookie that contains the refresh token
	refreshCookie, err := CreateRefreshCookie(guid, exp)
	if err != nil {
		return err
	}

	log.Printf("Created tokens for: %v\n", guid)
	// log.Println(accessCookie.Value)
	// log.Println(refreshCookie.Value)

	// Setting the cookies
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	return nil
}
