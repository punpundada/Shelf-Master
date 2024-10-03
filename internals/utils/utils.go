package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/punpundada/shelfMaster/internals/config"
	"github.com/punpundada/shelfMaster/internals/constants"
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

func WriteResponse(w http.ResponseWriter, code int, message string, result any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": message,
		"result":  result,
		"code":    code,
	})
}

func GetUserFromContext(cxt context.Context) (*db.User, error) {
	user, ok := cxt.Value(constants.User).(*db.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func GetSessionFromContext(ctx context.Context) (*db.Session, error) {
	session, ok := ctx.Value(constants.Session).(*db.Session)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func SendVerificationEmail(email string, code string) error {
	auth := smtp.PlainAuth("", config.GetConfig().SMTP_USERNAME, config.GetConfig().SMTP_PASSWORD, config.GetConfig().SMTP_HOST)
	from := config.GetConfig().SMTP_EMAIL
	to := []string{email}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = email
	headers["Subject"] = "Email Verification"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<p>Here's the OTP for your email verification:</p>
				<h2>%s</h2>
			</body>
		</html>`, code)

	var message string
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	smtpUrl := config.GetConfig().SMTP_HOST + ":" + config.GetConfig().SMTP_PORT
	err := smtp.SendMail(smtpUrl, auth, from, to, []byte(message))
	return err
}

func GenerateRandomDigits(n int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := strings.Builder{}
	for i := 0; i < n; i++ {
		digits.WriteString(strconv.Itoa(rnd.Intn(10)))
	}
	return digits.String()
}

func IsStrongPassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "password must contain at least 6 characters"
	}
	uppercase, _ := regexp.MatchString(`[A-Z]`, password)
	if !uppercase {
		return false, "password must contain at least one uppercase letter"
	}
	lowercase, _ := regexp.MatchString(`[a-z]`, password)
	if !lowercase {
		return false, "password must contain at least one lowercase letter"
	}
	digit, _ := regexp.MatchString(`[0-9]`, password)
	if !digit {
		return false, "password must contain at least one digit"
	}
	specialChar, _ := regexp.MatchString(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?~]`, password)
	if !specialChar {
		return false, "password must contain at least one special character"
	}
	return true, ""
}

func VerifyVerificationCode(ctx context.Context, db pgx.Tx, queries *db.Queries, user *db.User, code string) (bool, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)
	dbCode, err := qtx.GetEmailVerificationByUserId(ctx, user.ID)
	if err != nil {
		tx.Commit(ctx)
		return false, fmt.Errorf("user not found")
	}
	_, err = qtx.DeleteEmailVerificationByUserId(ctx, user.ID)
	if err != nil {
		tx.Rollback(ctx)
		return false, fmt.Errorf("code was not deleted")
	}
	tx.Commit(ctx)
	isNotExpired := isWithinExpirationDate(dbCode.ExpiresAt.Time)
	if !isNotExpired {
		return false, nil
	}
	if dbCode.Email != user.Email {
		return false, nil
	}
	return true, nil
}

func isWithinExpirationDate(expirationDate time.Time) bool {
	currentTime := time.Now()
	return expirationDate.After(currentTime)
}

func ParseJSON(request *http.Request, body any) error {
	err := json.NewDecoder(request.Body).Decode(body)
	defer request.Body.Close()
	return err
}
