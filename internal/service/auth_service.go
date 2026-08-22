package service

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	appjwt "go-waste-routes/internal/platform/jwt"
	"go-waste-routes/internal/domain"
)

type LoginResult struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         domain.User  `json:"user"`
}

type AuthService struct {
	users      *UserService
	jwtManager *appjwt.Manager
	sessions   map[string]int64
	failed     map[string]*failedLogin
	mu         sync.RWMutex
	lockout    time.Duration
	maxAttempts int
}

type failedLogin struct {
	Count       int
	LockedUntil time.Time
}

func NewAuthService(users *UserService, jwtManager *appjwt.Manager, maxAttempts int, lockout time.Duration) *AuthService {
	return &AuthService{
		users: users,
		jwtManager: jwtManager,
		sessions: make(map[string]int64),
		failed: make(map[string]*failedLogin),
		maxAttempts: maxAttempts,
		lockout: lockout,
	}
}

func (s *AuthService) Register(ctx context.Context, user domain.User, rawPassword string) (domain.User, error) {
	return s.users.Create(ctx, user, rawPassword, []domain.Permission{DefaultPermissions[2], DefaultPermissions[3]})
}

func (s *AuthService) Login(ctx context.Context, username, password, ip string) (LoginResult, error) {
	if err := s.checkLockout(ip); err != nil {
		return LoginResult{}, err
	}
	user, ok, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok || user.Status != domain.UserStatusEnabled {
		s.recordFailure(ip)
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	if err := s.users.VerifyPassword(user, password); err != nil {
		s.recordFailure(ip)
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	s.clearFailures(ip)
	access, refresh, tokenID, err := s.jwtManager.Generate(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, err
	}
	s.mu.Lock()
	s.sessions[tokenID] = user.ID
	s.mu.Unlock()
	return LoginResult{AccessToken: access, RefreshToken: refresh, User: user}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	claims, err := s.jwtManager.Parse(refreshToken, "refresh")
	if err != nil {
		return LoginResult{}, err
	}
	s.mu.RLock()
	_, ok := s.sessions[claims.ID]
	s.mu.RUnlock()
	if !ok {
		return LoginResult{}, fmt.Errorf("session revoked")
	}
	user, ok, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok {
		return LoginResult{}, fmt.Errorf("user not found")
	}
	access, refresh, tokenID, err := s.jwtManager.Generate(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, err
	}
	s.mu.Lock()
	delete(s.sessions, claims.ID)
	s.sessions[tokenID] = user.ID
	s.mu.Unlock()
	return LoginResult{AccessToken: access, RefreshToken: refresh, User: user}, nil
}

func (s *AuthService) Logout(refreshToken string) {
	claims, err := s.jwtManager.Parse(refreshToken, "refresh")
	if err != nil {
		return
	}
	s.mu.Lock()
	delete(s.sessions, claims.ID)
	s.mu.Unlock()
}

func (s *AuthService) Me(ctx context.Context, userID int64) (domain.User, error) {
	user, ok, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *AuthService) CheckAccess(token string) (*appjwt.Claims, error) {
	return s.jwtManager.Parse(token, "access")
}

func (s *AuthService) AuthenticatePassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *AuthService) checkLockout(ip string) error {
	s.mu.RLock()
	state := s.failed[ip]
	s.mu.RUnlock()
	if state == nil {
		return nil
	}
	if time.Now().Before(state.LockedUntil) {
		return fmt.Errorf("login locked until %s", state.LockedUntil.Format(time.RFC3339))
	}
	return nil
}

func (s *AuthService) recordFailure(ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.failed[ip]
	if state == nil {
		state = &failedLogin{}
		s.failed[ip] = state
	}
	state.Count++
	if state.Count >= s.maxAttempts {
		state.LockedUntil = time.Now().Add(s.lockout)
		state.Count = 0
	}
}

func (s *AuthService) clearFailures(ip string) {
	s.mu.Lock()
	delete(s.failed, ip)
	s.mu.Unlock()
}

func RemoteIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	if host == "" {
		return remoteAddr
	}
	return strings.TrimSpace(host)
}
