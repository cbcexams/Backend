package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Session struct {
	Id        string    `orm:"pk;size(36)" json:"id"`
	UserId    int       `orm:"column(user_id)" json:"user_id"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)" json:"created_at"`
	ExpiresAt time.Time `orm:"type(timestamp)" json:"expires_at"`
}

func (s *Session) TableName() string {
	return "session"
}

// CreateSession creates a new session and returns the session ID
func CreateSession() string {
	// Generate random session ID
	bytes := make([]byte, 32)
	rand.Read(bytes)
	sessionID := hex.EncodeToString(bytes)

	// Create session with 24 hour expiry
	session := &Session{
		Id:        sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	o := orm.NewOrm()
	_, err := o.Insert(session)
	if err != nil {
		return ""
	}

	return sessionID
}

// ValidateSession checks if a session is valid and not expired
func ValidateSession(sessionID string) error {
	o := orm.NewOrm()
	session := Session{Id: sessionID}

	if err := o.Read(&session); err != nil {
		return err
	}

	if time.Now().After(session.ExpiresAt) {
		o.Delete(&session)
		return orm.ErrNoRows
	}

	// Extend session expiry
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	_, err := o.Update(&session, "ExpiresAt")
	return err
}

// LinkSessionToUser associates a session with a user
func LinkSessionToUser(sessionID string, userID int) error {
	o := orm.NewOrm()
	session := Session{Id: sessionID}

	if err := o.Read(&session); err != nil {
		return err
	}

	session.UserId = userID
	_, err := o.Update(&session, "UserId")
	return err
}
