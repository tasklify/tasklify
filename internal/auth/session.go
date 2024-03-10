package auth

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"sync"
	"tasklify/internal/config"
	"tasklify/internal/database"
	"time"

	"github.com/wader/gormstore/v2"
)

const (
	cookieName      = "session"
	userIDFieldName = "user_id"
)

type Session interface {
	Create(userID uint, w http.ResponseWriter, r *http.Request) error
	Destroy(w http.ResponseWriter, r *http.Request) error
	GetUserID(r *http.Request) (uint, error)
}

type session struct {
	*gormstore.Store
}

var (
	onceSession sync.Once

	sessionClient *session
)

func GetSession(config ...*config.Config) Session {

	onceSession.Do(func() { // <-- atomic, does not allow repeating
		config := config[0]

		sessionClient = connectSession(config.Auth)

		// Periodic cleanup
		quit := make(chan struct{})                         // close quit channel to stop cleanup
		go sessionClient.PeriodicCleanup(1*time.Hour, quit) // db cleanup every hour
	})

	return sessionClient
}

func connectSession(config config.Auth) *session {

	sessionHashKey := sha512.Sum512([]byte(config.SessionHashKey))
	sessionBlockKey := sha256.Sum256([]byte(config.SessionBlockKey))

	store := gormstore.New(
		database.GetDatabase().RawDB(),
		sessionHashKey[:],
		sessionBlockKey[:],
	)

	store.SessionOpts.Secure = true
	store.SessionOpts.HttpOnly = true
	// store.SessionOpts.MaxAge = 2 * 24 * int(time.Hour)

	log.Println("Session connected")

	return &session{Store: store}
}

func (s *session) Create(userID uint, w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, cookieName)
	if err != nil {
		return err
	}

	session.Values[userIDFieldName] = userID
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) Destroy(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, cookieName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1 // Destroys cookie/session
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) GetUserID(r *http.Request) (uint, error) {
	session, err := s.Get(r, cookieName)
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values[userIDFieldName].(uint)
	if !ok {
		return 0, fmt.Errorf("could not get value of '%s'", userIDFieldName)
	}

	return userID, nil
}
