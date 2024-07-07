package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var errBadToken = errors.New("bad token")

var Validate = validator.New()

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenString, err := splitToken(tokenAuth); err == nil {
		return tokenString, nil
	}

	if tokenQuery != "" {
		return tokenQuery, nil
	}

	return "", errBadToken
}

func splitToken(token string) (string, error) {
	arr := strings.Split(token, " ")
	if len(arr) != 2 {
		return "", errBadToken
	}

	if arr[0] != "Bearer" {
		return "", errBadToken
	}

	return arr[1], nil
}

