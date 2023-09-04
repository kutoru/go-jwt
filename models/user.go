package models

import "time"

// Struct used for marshalling and unmarshalling
type User struct {
	Token string
	GUID  string
	EXP   time.Time
}
