package auth

import (
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
	userIdFieldName = "user_id"
)

type Session interface {
	GetUserId(r *http.Request) (string, error)
	SaveUserId(userId string, w http.ResponseWriter, r *http.Request) (http.ResponseWriter, error)
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
	ses := gormstore.New(
		database.GetDatabase().RawDB(),
		[]byte(config.SessionHashKey),
		[]byte(config.SessionBlockKey),
	)

	log.Println("Session connected")

	return &session{Store: ses}
}

func (s *session) GetUserId(r *http.Request) (string, error) {
	session, err := s.Get(r, cookieName)
	if err != nil {
		return "", err
	}

	userId, ok := session.Values[userIdFieldName].(string)
	if !ok {
		return "", fmt.Errorf("could not get value of '%s'", userIdFieldName)
	}

	return userId, nil
}

func (s *session) SaveUserId(userId string, w http.ResponseWriter, r *http.Request) (http.ResponseWriter, error) {
	session, err := s.Get(r, cookieName)
	if err != nil {
		return nil, err
	}

	session.Values[userIdFieldName] = userId
	sessionClient.Save(r, w, session)

	return w, nil
}
