package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/kutoru/go-jwt/db"
	"github.com/kutoru/go-jwt/glb"
	"github.com/kutoru/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Creates a cookie with a new refresh token as a value and inserts the hashed token into the DB
func CreateRefreshCookie(guid int, exp time.Time) (*http.Cookie, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	hashedTokenBytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedToken := string(hashedTokenBytes)
	err = insertRefreshTokenIntoDB(guid, hashedToken, exp)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "refresh",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   glb.EXP,
		SameSite: http.SameSiteNoneMode,
	}

	// log.Println("Created refresh token:", token)
	// log.Println("Created hashed refresh token:", hashedToken)

	return cookie, nil
}

// Fetches the refresh token info from the DB and returns it if the token is valid. Otherwise returns an error
func GetRefreshTokenInfo(r *http.Request, guid int) (int, int64, error) {
	cookie, err := r.Cookie("refresh")
	if err != nil {
		return 0, 0, err
	}

	collection := db.GetTokenCollection()
	filter := bson.M{"guid": guid}
	result := collection.FindOne(db.CTX, filter)

	var user models.User
	err = result.Decode(&user)
	if err != nil {
		return 0, 0, err
	}

	exp := user.EXP.Unix()
	if exp <= time.Now().Unix() {
		return 0, 0, fmt.Errorf("token has expired")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Token), []byte(cookie.Value),
	)
	if err != nil {
		return 0, 0, err
	}

	return user.GUID, exp, nil
}

// Inserts the refresh token into the DB
func insertRefreshTokenIntoDB(guid int, hashedToken string, exp time.Time) error {
	document := bson.M{
		"token": hashedToken,
		"guid":  guid,
		"exp":   exp,
	}

	collection := db.GetTokenCollection()
	_, err := collection.InsertOne(db.CTX, document)
	return err
}

// Removes document that has the guid
func RemoveRefreshTokenFromDB(guid int) error {
	collection := db.GetTokenCollection()
	filter := bson.M{"guid": guid}
	_, err := collection.DeleteOne(db.CTX, filter)
	return err
}
