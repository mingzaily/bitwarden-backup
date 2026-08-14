package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	sessionCookieName = "bw_session"
	sessionTTL        = 12 * time.Hour
	maxLoginAttempts  = 10
	loginWindow       = 15 * time.Minute
	maxActiveSessions = 256
	maxTrackedClients = 4096
)

type session struct {
	csrfToken string
	csrfHash  [32]byte
	expiresAt time.Time
}

type loginAttempt struct {
	count   int
	resetAt time.Time
}

// Manager provides a small, in-memory session layer for the single-admin UI.
// Sessions intentionally disappear when the process restarts, so a restart also
// invalidates every browser session.
type Manager struct {
	passwordHash [32]byte
	cookieSecure bool

	mu       sync.Mutex
	sessions map[[32]byte]session
	attempts map[string]loginAttempt
}

func New(password string, cookieSecure bool) *Manager {
	return &Manager{
		passwordHash: sha256.Sum256([]byte(password)),
		cookieSecure: cookieSecure,
		sessions:     make(map[[32]byte]session),
		attempts:     make(map[string]loginAttempt),
	}
}

// Login authenticates the administrator and issues an HttpOnly session cookie.
func (m *Manager) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	clientKey := clientAddress(c.Request)
	if !m.allowLogin(clientKey) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts, try again later"})
		return
	}

	providedHash := sha256.Sum256([]byte(req.Password))
	if subtle.ConstantTimeCompare(providedHash[:], m.passwordHash[:]) != 1 {
		m.recordFailedLogin(clientKey)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	m.clearLoginAttempts(clientKey)
	token, err := randomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}
	csrfToken, err := randomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	tokenHash := sha256.Sum256([]byte(token))
	csrfHash := sha256.Sum256([]byte(csrfToken))
	m.mu.Lock()
	m.purgeExpiredLocked(time.Now())
	if len(m.sessions) >= maxActiveSessions {
		m.removeOldestSessionLocked()
	}
	m.sessions[tokenHash] = session{
		csrfToken: csrfToken,
		csrfHash:  csrfHash,
		expiresAt: time.Now().Add(sessionTTL),
	}
	m.mu.Unlock()

	m.setSessionCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"authenticated": true, "csrf_token": csrfToken})
}

// Require protects API routes with the session cookie.
func (m *Manager) Require() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := m.sessionForRequest(c); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		c.Next()
	}
}

// CSRF protects state-changing requests. SameSite cookies are the first line
// of defense; this header check also blocks cross-site form and fetch attempts.
func (m *Manager) CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		sess, ok := m.sessionForRequest(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		provided := c.GetHeader("X-CSRF-Token")
		providedHash := sha256.Sum256([]byte(provided))
		if provided == "" || subtle.ConstantTimeCompare(providedHash[:], sess.csrfHash[:]) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			return
		}
		c.Next()
	}
}

// Session returns the current session's CSRF token to the SPA. The session
// cookie remains HttpOnly and is never exposed to JavaScript.
func (m *Manager) Session(c *gin.Context) {
	sess, ok := m.sessionForRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authenticated": true, "csrf_token": sess.csrfToken})
}

func (m *Manager) Logout(c *gin.Context) {
	if token, err := c.Cookie(sessionCookieName); err == nil && token != "" {
		tokenHash := sha256.Sum256([]byte(token))
		m.mu.Lock()
		delete(m.sessions, tokenHash)
		m.mu.Unlock()
	}
	m.clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"authenticated": false})
}

func (m *Manager) sessionForRequest(c *gin.Context) (session, bool) {
	encoded, err := c.Cookie(sessionCookieName)
	if err != nil || encoded == "" {
		return session{}, false
	}
	tokenHash := sha256.Sum256([]byte(encoded))

	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[tokenHash]
	if !ok {
		return session{}, false
	}
	if time.Now().After(sess.expiresAt) {
		delete(m.sessions, tokenHash)
		return session{}, false
	}
	return sess, true
}

func (m *Manager) setSessionCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(sessionCookieName, token, int(sessionTTL.Seconds()), "/", "", m.cookieSecure, true)
}

func (m *Manager) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(sessionCookieName, "", -1, "/", "", m.cookieSecure, true)
}

func (m *Manager) allowLogin(clientKey string) bool {
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.purgeExpiredLocked(now)
	attempt, ok := m.attempts[clientKey]
	if !ok || now.After(attempt.resetAt) {
		return true
	}
	return attempt.count < maxLoginAttempts
}

func (m *Manager) recordFailedLogin(clientKey string) {
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.purgeExpiredLocked(now)
	if _, exists := m.attempts[clientKey]; !exists && len(m.attempts) >= maxTrackedClients {
		return
	}
	attempt := m.attempts[clientKey]
	if attempt.resetAt.IsZero() || now.After(attempt.resetAt) {
		attempt = loginAttempt{resetAt: now.Add(loginWindow)}
	}
	attempt.count++
	m.attempts[clientKey] = attempt
}

func (m *Manager) clearLoginAttempts(clientKey string) {
	m.mu.Lock()
	delete(m.attempts, clientKey)
	m.mu.Unlock()
}

func (m *Manager) purgeExpiredLocked(now time.Time) {
	for tokenHash, sess := range m.sessions {
		if now.After(sess.expiresAt) {
			delete(m.sessions, tokenHash)
		}
	}
	for clientKey, attempt := range m.attempts {
		if !attempt.resetAt.IsZero() && now.After(attempt.resetAt) {
			delete(m.attempts, clientKey)
		}
	}
}

func (m *Manager) removeOldestSessionLocked() {
	var oldestToken [32]byte
	var oldest time.Time
	for tokenHash, sess := range m.sessions {
		if oldest.IsZero() || sess.expiresAt.Before(oldest) {
			oldestToken = tokenHash
			oldest = sess.expiresAt
		}
	}
	if !oldest.IsZero() {
		delete(m.sessions, oldestToken)
	}
}

func randomToken() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func clientAddress(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	if strings.TrimSpace(r.RemoteAddr) != "" {
		return r.RemoteAddr
	}
	return "unknown"
}
