package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/punpundada/shelfMaster/internals/config"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type SuccessResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Result    any    `json:"result"`
}

func HashString(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyHashString(hashedString, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(str))
	return err == nil
}

func VerifyRequestOrigin(origin string, hosts []string) bool {
	parsedUrl, err := url.Parse(origin)
	if err != nil {
		return false
	}
	originHost := parsedUrl.Hostname()

	for _, host := range hosts {
		if host == origin || strings.HasPrefix(originHost, "."+host) {
			return true
		}
	}
	return false
}

func ValidateSession(ctx context.Context, queries *db.Queries, sessionId string) (*db.Session, *db.User, error) {
	session, err := queries.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, nil, err
	}
	user, err := queries.GetUserById(ctx, session.UserID)

	if err != nil {
		return nil, nil, err
	}
	if time.Now().After(session.ExpiresAt.Time) {
		session.Fresh = pgtype.Bool{Bool: false, Valid: true}
	}
	return &session, &user, nil
}

func CreateSessionCookies(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 14 * 24), //2 weeks
		Secure:   config.GetConfig().ENV != "development",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
}
func CreateBlankSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   config.GetConfig().ENV != "development",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
}

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string, details ...string) {
	errorResponse := ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		errorResponse.Details = details[0]
		var sb strings.Builder
		for index, item := range details {
			if index == 1 {
				sb.WriteString(item)
			}
			sb.WriteString(". " + item)
		}
	}

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Failed to generate error response", http.StatusInternalServerError)
	}
}

func GetUserFromContext(cxt context.Context) (*db.User, error) {
	user, ok := cxt.Value(db.RoleTypeUSER).(db.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func GetSessionFromContext(ctx context.Context) (*db.Session, error) {
	session, ok := ctx.Value("SESSION").(db.Session)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return &session, nil
}
