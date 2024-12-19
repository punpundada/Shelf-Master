package utils

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
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

	"github.com/google/uuid"
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

type ApiError struct {
	Message string
	Code    int
}

func (a *ApiError) Error() string {
	return a.Message
}

func (a *ApiError) WriteError(w http.ResponseWriter, details ...string) {
	errorResponse := ErrorResponse{
		Success: false,
		Code:    a.Code,
		Message: a.Message,
	}
	if len(details) > 0 {
		var sb strings.Builder
		for index, item := range details {
			if index == 1 {
				sb.WriteString(item)
			}
			sb.WriteString(". " + item)
		}
		errorResponse.Details = sb.String()
	}

	w.WriteHeader(a.Code)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Failed to generate error response", http.StatusInternalServerError)
	}
}

func NewApiError(msg string, code int) *ApiError {
	return &ApiError{
		Message: msg,
		Code:    code,
	}
}

func HashString(str string) (string, *ApiError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", NewApiError(err.Error(), http.StatusInternalServerError)
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

func ValidateSession(ctx context.Context, queries *db.Queries, sessionId string) (*db.Session, *db.User, *ApiError) {
	session, err := queries.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, nil, NewApiError(err.Error(), http.StatusBadRequest)
	}
	user, err := queries.GetUserById(ctx, session.UserID)

	if err != nil {
		return nil, nil, NewApiError(err.Error(), http.StatusBadRequest)
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

func GetUserFromContext(cxt context.Context) (*db.User, *ApiError) {
	user, ok := cxt.Value(constants.User).(*db.User)
	if !ok {
		return nil, NewApiError("user not found", http.StatusForbidden)
	}
	return user, nil
}

func GetSessionFromContext(ctx context.Context) (*db.Session, *ApiError) {
	session, ok := ctx.Value(constants.Session).(*db.Session)
	if !ok {
		return nil, NewApiError("session not found", http.StatusForbidden)
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

func VerifyVerificationCode(ctx context.Context, db pgx.Tx, queries *db.Queries, user *db.User, code string) (bool, *ApiError) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return false, NewApiError(err.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)
	dbCode, err := qtx.GetEmailVerificationByUserId(ctx, user.ID)
	if err != nil {
		tx.Commit(ctx)
		return false, NewApiError("user not found", http.StatusBadRequest)
	}
	_, err = qtx.DeleteEmailVerificationByUserId(ctx, user.ID)
	if err != nil {
		tx.Rollback(ctx)
		return false, NewApiError("user not found", http.StatusBadRequest)
	}
	tx.Commit(ctx)
	isvalid := IsWithinExpirationDate(dbCode.ExpiresAt.Time)
	if !isvalid {
		return false, nil
	}
	if dbCode.Email != user.Email {
		return false, nil
	}
	return true, nil
}

func IsWithinExpirationDate(expirationDate time.Time) bool {
	currentTime := time.Now()
	return !expirationDate.After(currentTime)
}

func ParseJSON(request *http.Request, body any) *ApiError {
	err := json.NewDecoder(request.Body).Decode(body)
	defer request.Body.Close()
	if err != nil {
		return NewApiError("error parsing body "+err.Error(), http.StatusBadRequest)
	}
	return nil
}

func MarshalJson(w http.ResponseWriter, body any) error {
	return json.NewEncoder(w).Encode(&body)
}

func NewSaveSessionAttrs(userId int32) *db.SaveSessionParams {
	return &db.SaveSessionParams{
		ID:     uuid.New().String(),
		UserID: userId,
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Hour * 24 * 2),
			Valid: true,
		},
	}
}

func SendPasswordResetEmail(email string, code string) error {
	auth := smtp.PlainAuth("", config.GetConfig().SMTP_USERNAME, config.GetConfig().SMTP_PASSWORD, config.GetConfig().SMTP_HOST)
	from := config.GetConfig().SMTP_EMAIL
	to := []string{email}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = email
	headers["Subject"] = "Password Reset"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	htmlBody := fmt.Sprintf(`
	<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
					color: #333;
				}
				.container {
					max-width: 600px;
					margin: 50px auto;
					padding: 20px;
					background-color: #fff;
					border-radius: 8px;
					box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #444;
					font-size: 24px;
					text-align: center;
				}
				p {
					font-size: 16px;
					line-height: 1.6;
					color: #666;
					text-align: center;
				}
				.reset-button {
					display: inline-block;
					padding: 10px 20px;
					font-size: 16px;
					color: white;
					background-color: #007BFF;
					text-decoration: none;
					border-radius: 5px;
					text-align: center;
				}
				.reset-button:hover {
					background-color: #0056b3;
				}
				.footer {
					font-size: 12px;
					text-align: center;
					color: #888;
					margin-top: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Reset Your Password</h1>
				<p>You requested a password reset. Click the button below to reset your password:</p>
				<p><a href="%s" class="reset-button" style="color:white">Reset Password</a></p>
				<div class="footer">
					<p>If you did not request a password reset, please ignore this email.</p>
				</div>
			</div>
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

func CreatePasswordRestToken(ctx context.Context, q *db.Queries, userId int32) (string, error) {
	if err := q.DeleteRestPasswordByUserId(ctx, userId); err != nil {
		return "", err
	}
	tokenId, err := generateIdFromEntropy(25)
	if err != nil {
		return "", err
	}
	tokenHash := EncodeString(tokenId)
	_, err = q.SavePasswordRestToken(ctx, db.SavePasswordRestTokenParams{
		TokenHash: pgtype.Text{String: tokenHash, Valid: true},
		UserID:    userId,
		ExpiresAt: pgtype.Date{
			Time:  time.Now().Add(time.Minute * 15),
			Valid: true,
		},
	})
	if err != nil {
		return "", err
	}

	return tokenId, nil
}

func generateIdFromEntropy(size int) (string, error) {
	// Create a byte slice with the desired entropy size
	buffer := make([]byte, size)

	// Fill the slice with random values
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))
	_, err := rnd.Read(buffer)
	if err != nil {
		return "", err
	}

	// Encode the random values using Base32 encoding
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buffer)

	// Convert to lowercase as per your JS implementation
	encoded = strings.ToLower(encoded)

	return encoded, nil
}

func InvalidateAllUserSessions(ctx context.Context, q *db.Queries, userId int32) error {
	_, err := q.DeleteSessionByUserId(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func EncodeString(verificationToken string) string {
	tokenBytes := []byte(verificationToken)
	hash := sha256.Sum256(tokenBytes)
	hexHash := hex.EncodeToString(hash[:]) //hash[:] converts [32]byte into []byte i.e. array -> slice
	return hexHash
}
