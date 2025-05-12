package mkusersessionservice

import (
	"sync"

	"google.golang.org/api/oauth2/v2"

	"github.com/google/uuid"
)

type UserInfo struct {
	Id                int64
	GoogleCredentials []*oauth2.Userinfo
}

type SessionMap struct {
	m map[string]UserInfo
	l sync.RWMutex
}

type Service struct {
	sessionMap SessionMap
	nextUserId int64
}

func (s *Service) CreateUserSession(googleCredentials []*oauth2.Userinfo) (string, error) {
	s.sessionMap.l.Lock()
	sessionId := uuid.New().String()
	s.sessionMap.m[sessionId] = UserInfo{
		Id:                s.nextUserId,
		GoogleCredentials: googleCredentials,
	}

	s.nextUserId++
	s.sessionMap.l.Unlock()

	return sessionId, nil
}

func (s *Service) GetSession(id string) UserInfo {
	s.sessionMap.l.RLock()
	defer s.sessionMap.l.RUnlock()
	return s.sessionMap.m[id]
}

func New() *Service {
	return &Service{
		sessionMap: SessionMap{
			m: make(map[string]UserInfo),
			l: sync.RWMutex{},
		},
		nextUserId: 0,
	}
}
