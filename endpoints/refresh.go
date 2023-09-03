package endpoints

import (
	"log"
	"net/http"

	"github.com/kutoru/go-jwt/tokens"
)

// Function that handles the refresh route
func refresh(w http.ResponseWriter, r *http.Request) {
	log.Println("\nrefresh called")

	// Getting both tokens' info and checking if they are valid

	accessGuid, accessExp, err := tokens.GetAccessTokenInfo(r)
	if err != nil {
		log.Printf("Invalid access token: %v\n", err)
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	refreshGuid, refreshExp, err := tokens.GetRefreshTokenInfo(r, accessGuid)
	if err != nil {
		log.Printf("Invalid refresh token: %v\n", err)
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	// Comparing the info
	if accessGuid != refreshGuid || accessExp != refreshExp {
		log.Printf("The tokens do not match")
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	// Removing the old refresh token
	err = tokens.RemoveRefreshTokenFromDB(accessGuid)
	if err != nil {
		log.Printf("Could not delete the old refresh token")
		http.Error(w, "Could not update token", http.StatusInternalServerError)
		return
	}

	err = tokens.CreateAndSetAsCookies(w, accessGuid)
	if err != nil {
		log.Printf("Could not create tokens: %v\n", err)
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(make([]byte, 0))
	if err != nil {
		log.Printf("Could not write a response: %v\n", err)
	}
}
