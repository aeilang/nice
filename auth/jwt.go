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
	Email string `json:"email"`
	Role  string `json:"role"`
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

			payload, err := ValidateJWT(tokenString, configs.Envs.JWTAccessSecret)
			if err != nil {
				permissionDenied(w)
				log.Println(err)
				return
			}

			u, err := querier.GetUserByEmail(context.Background(), payload.Email)
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				permissionDenied(w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, userKey, u)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func CreateJWT(secret string, email string, role string, expiration time.Duration) (string, error) {

	claims := JwtPayload{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, secret string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
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

func GetUserFromContext(ctx context.Context) (store.User, bool) {
	user, ok := ctx.Value(userKey).(store.User)
	return user, ok
}
