package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/Ashutoshbind15/gogameengine/internal/data"
)

func GenerateSessionId() string {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	sessionId := base32.StdEncoding.EncodeToString(bytes)
	return sessionId
}

// todo: call the data-access layer later on
func CreateDBSession(uid int, sid string) {

	validTill := time.Now().Add(1 * time.Hour)
	dberr := data.DbConn.Ping();

	if dberr != nil {
		fmt.Println("db isnt connected yet: ", dberr)
		return
	}

	_, err := data.DbConn.Exec(`INSERT INTO DBSESSIONS (id, user_id, validTo) VALUES ($1, $2, $3)`,  sid, uid, validTill)
	if err != nil {
		fmt.Println("Error creating a sess in the db: ", err)
	}
}