package service

import (
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/jwtutil"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo     *repository.UserRepo
	onlineRepo repository.OnlineRepo
}

func NewUserService(repo *repository.UserRepo, onlineRepo repository.OnlineRepo) *UserService {
	return &UserService{repo: repo, onlineRepo: onlineRepo}
}

func (s *UserService) Register(username, email, password string) (*model.User, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetByEmail(email); err == nil {
		return nil, errors.New("email already registered")
	}
	if _, err := s.repo.GetByUsername(username); err == nil {
		return nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username, Email: email,
		Password: string(hash), GroupID: 101,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) Login(account, password, ip string) (*model.User, string, error) {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	now := time.Now()
	s.repo.Update(user.ID, map[string]interface{}{
		"last_login": now, "last_ip": ip,
	})
	user.LastLogin = &now

	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Track online user (non-blocking)
	_ = s.onlineRepo.Upsert(user.ID, user.Username, ip)

	return user, token, nil
}

func (s *UserService) GetByID(id uint64) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(id uint64, updates map[string]interface{}) error {
	return s.repo.Update(id, updates)
}

func (s *UserService) ChangePassword(id uint64, oldPwd, newPwd string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return errors.New("old password incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Update(id, map[string]interface{}{"password": string(hash)})
}

func (s *UserService) Delete(id uint64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteWithCleanup(id, user.GroupID)
}

func (s *UserService) List(page, pageSize int, search string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, search)
}

func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *UserService) GetGroup(db *gorm.DB, gid uint) (*model.Group, error) {
	var group model.Group
	err := db.First(&group, gid).Error
	return &group, err
}

func (s *UserService) TrackOnline(userID uint64, username, ip string) error {
	return s.onlineRepo.Upsert(userID, username, ip)
}

func (s *UserService) ListOnlineUsers(limit int) ([]model.OnlineUserWithInfo, error) {
	return s.onlineRepo.List(limit)
}

func (s *UserService) ListGroups() ([]model.Group, error) {
	return s.repo.ListGroups()
}

func (s *UserService) UpdateGroup(id uint, updates map[string]interface{}) error {
	return s.repo.UpdateGroup(id, updates)
}

func (s *UserService) VerifyEmail(id uint64) error {
	return s.repo.VerifyEmail(id)
}

func generateToken(user *model.User) (string, error) {
	return jwtutil.GenerateToken(user.ID, user.GroupID)
}

// GetIPFromRequest extracts real IP from request.
// Only trusts X-Forwarded-For/X-Real-IP headers if the remote IP is in trustedProxies.
// For Gin handlers, prefer c.ClientIP() with SetTrustedProxies configured.
func GetIPFromRequest(r *http.Request, trustedProxies []string) string {
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if remoteIP == "" {
		remoteIP = r.RemoteAddr
	}

	if !isTrustedProxy(remoteIP, trustedProxies) {
		return remoteIP
	}

	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		for _, p := range parts {
			ip := strings.TrimSpace(p)
			if ip != "" {
				return ip
			}
		}
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}
	return remoteIP
}

func isTrustedProxy(ip string, trustedProxies []string) bool {
	if len(trustedProxies) == 0 {
		return false
	}
	for _, p := range trustedProxies {
		if p == ip || p == "*" {
			return true
		}
	}
	return false
}

// ValidateUsername checks username validity
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 24 {
		return errors.New("username must be at most 24 characters")
	}
	// Allow alphanumeric, underscore, hyphen, Chinese characters
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_\-\x{4e00}-\x{9fa5}]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscore, hyphen, and Chinese characters")
	}
	// Reserved names
	reserved := []string{"admin", "root", "system", "moderator", "moderators", "guest", "null", "undefined"}
	for _, r := range reserved {
		if username == r {
			return errors.New("this username is reserved")
		}
	}
	return nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return errors.New("password must be at most 128 characters")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase letters and numbers")
	}
	return nil
}

// ValidateEmail checks email validity
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	if len(email) > 128 {
		return errors.New("email is too long")
	}
	return nil
}
