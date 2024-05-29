package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aeilang/nice/configs"
	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/utils"
	"github.com/golang-jwt/jwt/v5"
)

var userKey = struct{}{}

type JwtPayload struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func WithJWTAuth(querier store.Querier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.GetTokenFromRequest(r)
			if err != nil {
				permissionDenied(w)
				log.Println(err)
				return
			}

			payload, err := validateJWT(tokenString)
			if err != nil {
				permissionDenied(w)
				log.Println(err)
				return
			}

			u, err := querier.GetUserByID(r.Context(), int32(payload.UserID))
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				permissionDenied(w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, userKey, u.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Duration(configs.Envs.JWTExperationInMinites*60) * time.Second

	claims := JwtPayload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lang",
			Subject:   "lang",
			ID:        "1",
			Audience:  []string{"only myself"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*JwtPayload)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return payload, nil
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userKey).(int)
	return userID, ok
}
