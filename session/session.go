package session

import (
	"strconv"
	"errors"
	"math/rand"
	"sync"
	"net/http"
	"networks2016-task1/users"
	"time"
)

var Manager *SessionManager

func checkToken(token string) bool {
	return true // TODO: it's temporary
}

type SessionInterface interface {
	StartSession()
	GetSession()
	DeleteSession()
	SessionGC()
}

type Session struct {
	User	*users.User
}

func NewSession(user *users.User) (*Session) {
	return &Session{
		User:user,
	}
}

type SessionManager struct {
	ManagerName		string
	sessions		map[int]*Session
	randomGenerator	*rand.Rand
	mutex 			sync.Mutex
}

func NewSessionManager(managerName string) (*SessionManager, error) {
	return &SessionManager{
		ManagerName: managerName,
		sessions: make(map[int]*Session),
		randomGenerator: rand.New(rand.NewSource(time.Now().Unix())), // TODO: change seed to something random
	}, nil
}

func (manager *SessionManager) getTokenFromRequest(request *http.Request) (string, error) {
	cookie, err := request.Cookie(manager.ManagerName)
	if checkToken(cookie.Value) {
		return cookie.Value, err
	}

	return "", errors.New("Wrong token writen in cookie")
}

func (manager *SessionManager) generateNewToken() (int) {
	for {
		tokenSession := manager.randomGenerator.Int()
		_, ok := manager.sessions[tokenSession]
		if !ok {
			return tokenSession
		}
	}
}

func (manager *SessionManager) StartSession(responseWriter http.ResponseWriter, request *http.Request, user *users.User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	token := manager.generateNewToken()
	manager.sessions[token] = NewSession(user)

	cookie := http.Cookie{Name:manager.ManagerName, Value:strconv.Itoa(token), Path:"/", HttpOnly:true}
	http.SetCookie(responseWriter, &cookie)
	http.Redirect(responseWriter, request, "/", 302) // TODO: может быть редирект нужно делать не здесь
}

func (manager *SessionManager) GetSession(request *http.Request) (*Session, error) {
	token, err := manager.getTokenFromRequest(request)
	if err != nil {
		return nil, err
	}

	tokenSession, _ := strconv.Atoi(token)

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	session, ok := manager.sessions[tokenSession]
	if !ok {
		return nil, errors.New("ERROR: session " + strconv.Itoa(tokenSession) + " doesn't exists")
	}

	return session, nil
}

func (manager *SessionManager) IsLogged(responseWriter http.ResponseWriter, request *http.Request) (bool) {
	cookie, err := request.Cookie(manager.ManagerName)
	if err != nil || cookie.Value == "" {
		return false
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	token, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return false
	}

	_, ok := manager.sessions[token]
	return ok
}

func (manager *SessionManager) Logout(responseWriter http.ResponseWriter, request *http.Request) (error) {
	token, err := manager.getTokenFromRequest(request)
	if err != nil {
		return err
	}

	tokenSession, _ := strconv.Atoi(token)

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.sessions, tokenSession)
	return nil
}
