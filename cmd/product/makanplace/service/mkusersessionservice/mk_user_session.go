package mkusersessionservice

import (
	"sync"

	"golang-server/cmd/product/makanplace/repositories/auth"
	"google.golang.org/api/oauth2/v2"

	"github.com/google/uuid"
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
	sessionMap SessionMap
	authRepo   *auth.Repo
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

	s.sessionMap.m[sessionId] = user

	return sessionId, nil
}

// GetSession if hardened, dont return confidential values
func (s *Service) GetSession(id string, hard bool) *auth.UserWithGmail {
	s.sessionMap.l.RLock()
	defer s.sessionMap.l.RUnlock()
	session, ok := s.sessionMap.m[id]
	if !ok {
		return nil
	}
	if hard {
		session.Gmails = nil
	}
	return session
}

func (s *Service) RemoveSession(sessionId string) error {
	s.sessionMap.l.Lock()
	defer s.sessionMap.l.Unlock()
	delete(s.sessionMap.m, sessionId)

	return nil
}

func New(authRepo *auth.Repo) *Service {
	return &Service{
		sessionMap: SessionMap{
			m: make(map[string]*auth.UserWithGmail),
			l: sync.RWMutex{},
		},
		authRepo: authRepo,
	}
}
