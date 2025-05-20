package mk_user_session

import (
	"log"
	"slices"
	"sync"

	"golang-server/cmd/product/makanplace/config"
	"golang-server/cmd/product/makanplace/repositories/auth"

	"github.com/google/uuid"
	"google.golang.org/api/oauth2/v2"
)

type UserInfo struct {
	Id                int64              `json:"Id"`
	GoogleCredentials []*oauth2.Userinfo `json:"-"`
}

type SessionMap struct {
	m map[string]*auth.UserWithGmail
	l sync.RWMutex
}

type Service struct {
	sessionMap  SessionMap
	authRepo    *auth.Repo
	superGmails []string
}

func (s *Service) IsSuperUser(gmail string) bool {
	if len(s.superGmails) == 0 {
		return false
	}
	return slices.Contains(s.superGmails, gmail)
}

func ToEmails(a []*oauth2.Userinfo) (b []string) {
	for _, v := range a {
		b = append(b, v.Email)
	}
	return
}

func (s *Service) CreateUserSession(googleCredentials []*oauth2.Userinfo) (string, error) {
	s.sessionMap.l.Lock()
	defer s.sessionMap.l.Unlock()
	emails := ToEmails(googleCredentials)
	user, err := s.authRepo.GetOrCreateUserByGmail(emails)
	if err != nil {
		return "", err
	}

	sessionId := uuid.New().String()
	log.Printf("[CreateUserSession] sessionId: %#v, user %#v\n", sessionId, user)
	s.sessionMap.m[sessionId] = user
	return sessionId, nil
}

// GetSession if hardened, dont return confidential values
func (s *Service) GetSession(sessionId string, hard bool) *auth.UserWithGmail {
	if sessionId == "" {
		return nil
	}
	s.sessionMap.l.RLock()
	defer s.sessionMap.l.RUnlock()
	session, ok := s.sessionMap.m[sessionId]
	if !ok {
		return nil
	}
	if hard {
		_session := *session
		_session.Gmails = nil
		return &_session
	}
	return session
}

func (s *Service) RemoveSession(sessionId string) error {
	s.sessionMap.l.Lock()
	defer s.sessionMap.l.Unlock()
	delete(s.sessionMap.m, sessionId)

	return nil
}

func New(authRepo *auth.Repo, config config.AdminConfig) *Service {
	return &Service{
		sessionMap: SessionMap{
			m: make(map[string]*auth.UserWithGmail),
			l: sync.RWMutex{},
		},
		authRepo:    authRepo,
		superGmails: config.Gmails,
	}
}
