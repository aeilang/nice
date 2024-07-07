package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aeilang/nice/auth"
	"github.com/aeilang/nice/configs"
	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errs))
		return
	}

	u, err := s.Querier.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(u.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	// check if user exists
	_, err := s.Querier.GetUserByEmail(r.Context(), user.Email)
	if err == nil {
		log.Println(err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	u, err := s.Querier.CreateUser(r.Context(), store.CreateUserParams{
		Name:     user.UserName,
		Email:    user.Email,
		Password: hashedPassword,
		Role:     "user",
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, RegisterResponse{
		Username: u.Name,
		Email:    u.Email,
	})
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing path parameter"))
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("uncorrect path parameter"))
		return
	}

	u, err := s.Querier.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, GetUserResponse{
		Id:        u.ID,
		Username:  u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Role:      string(u.Role),
	})
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserPayload struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetUserResponse struct {
	Id        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"eamil"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}
