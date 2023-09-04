package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kutoru/go-jwt/models"
	"github.com/kutoru/go-jwt/tokens"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function that handles the auth route
func auth(w http.ResponseWriter, r *http.Request) {
	log.Println("\nauth called")

	// Getting GUID from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Could not parse the body content: %v\n", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Creating tokens and setting them as cookies to the writer
	err = tokens.CreateAndSetAsCookies(w, user.GUID)
	switch err.(type) {
	case mongo.WriteException:
		log.Printf("Could not create refresh token: %v\n", err)
		http.Error(w, "GUID already exists", http.StatusBadRequest)
		return
	default:
		log.Printf("Could not create tokens: %v\n", err)
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	case nil:
	}

	// Writing the response
	_, err = w.Write(make([]byte, 0))
	if err != nil {
		log.Printf("Could not write a response: %v\n", err)
	}
}
