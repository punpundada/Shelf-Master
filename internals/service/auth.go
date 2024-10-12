package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type AuthService struct {
	Queries *db.Queries
}

type LoginBody struct {
	Email string `json:"email"`
}

func (a *AuthService) LoginUser(r *http.Request) (*db.User, *db.Session, error) {
	var body = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	if (len(body.Email) == 0) || (!utils.IsValidEmail(body.Email)) {
		return nil, nil, fmt.Errorf("invalid email")
	}
	user, err := a.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found either password or email do not match: %v", err)
	}
	isCorrectPassword := utils.VerifyHashString(user.PasswordHash, body.Password)
	if !isCorrectPassword {
		return nil, nil, fmt.Errorf("user not found either password or email do not match")
	}
	fmt.Println(user)
	sessionId := uuid.New()
	session, err := a.Queries.SaveSession(r.Context(), db.SaveSessionParams{
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 14)},
		ID:        sessionId.String(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while saving session: %v", err)
	}
	return &user, &session, nil
}

func (a *AuthService) SaveUser(r *http.Request) (*db.User, error) {
	var body db.SaveUserParams
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error decoding body: %v", err)
	}
	defer r.Body.Close()
	if isValidEmail := utils.IsValidEmail(body.Email); !isValidEmail {
		return nil, fmt.Errorf("invalid email")
	}
	if isStrongPassword, msg := utils.IsStrongPassword(body.PasswordHash); !isStrongPassword {
		return nil, fmt.Errorf("inscure password: %s", msg)
	}
	hashedPassword, err := utils.HashString(body.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}
	body.PasswordHash = hashedPassword

	user, err := a.Queries.SaveUser(r.Context(), body)
	if err != nil {
		if strings.Contains(err.Error(), "email_unique") {
			user, err := a.Queries.GetUserByEmail(r.Context(), body.Email)
			if err != nil {
				return nil, fmt.Errorf("no user found: %v", err)
			}
			if !user.EmailVerified.Bool {
				return nil, fmt.Errorf("email already in use")
			}
		}
		return nil, fmt.Errorf("error saving user: %v", err)
	}
	verificationCode, err := generateEmailVerificationCode(r.Context(), user.ID, user.Email, a.Queries)
	if err != nil {
		return nil, fmt.Errorf("error generating verification code or saving code: %v", err)
	}
	err = utils.SendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		return nil, fmt.Errorf("error sending email: %v", err)
	}
	return &user, nil
}

func generateEmailVerificationCode(ctx context.Context, userId int32, email string, q *db.Queries) (string, error) {
	_, err := q.DeleteEmailVerificationByUserId(ctx, userId)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return "", fmt.Errorf("error deleting verifications: %v", err)
		}
	}

	code := utils.GenerateRandomDigits(6)
	data := db.SaveEmailVerificationParams{
		Code:      code,
		UserID:    userId,
		Email:     email,
		ExpiresAt: pgtype.Date{Time: time.Now().Add(15 * time.Minute), Valid: true},
	}
	_, err = q.SaveEmailVerification(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error saving varification code: %v", err)
	}
	return code, nil
}
