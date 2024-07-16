package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aeilang/nice/auth"
	"github.com/aeilang/nice/configs"
	"github.com/aeilang/nice/internal/repository/store"
	"github.com/aeilang/nice/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("用户邮箱不存在"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("密码不正确"))
		return
	}

	accessExpiration := time.Duration(configs.Envs.JWTAccessExperationInMinites) * time.Minute
	accessToken, err := auth.CreateJWT(configs.Envs.JWTAccessSecret, u.Email, string(u.Role), accessExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshExpiration := time.Duration(configs.Envs.JWTRefreshExperationInHours) * time.Hour
	refreshToken, err := auth.CreateJWT(configs.Envs.JWTRefreshSecret, u.Email, string(u.Role), refreshExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Rdb.Set(context.Background(), u.Email, refreshToken, refreshExpiration).Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *Server) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var token RefreshTokenPayload
	if err := utils.ParseJSON(r, &token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payload, err := auth.ValidateJWT(token.RefreshToken, configs.Envs.JWTRefreshSecret)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("解析token错误"))
		return
	}

	refreshToken, err := s.Rdb.Get(context.Background(), payload.Email).Result()
	if err != nil || refreshToken != token.RefreshToken {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("refreshToken失效"))
		return
	}

	u, err := s.Querier.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("邮箱不存在"))
		return
	}

	accessExpiration := time.Duration(configs.Envs.JWTAccessExperationInMinites) * time.Minute
	accessToken, err := auth.CreateJWT(configs.Envs.JWTAccessSecret, u.Email, string(u.Role), accessExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshExpiration := time.Duration(configs.Envs.JWTRefreshExperationInHours) * time.Hour
	newRefreshToken, err := auth.CreateJWT(configs.Envs.JWTRefreshSecret, u.Email, string(u.Role), refreshExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Rdb.Set(context.Background(), u.Email, newRefreshToken, refreshExpiration).Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})

}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
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

	val, err := s.Rdb.Get(context.Background(), user.Email).Result()
	if err != nil || val != user.VerifiCode {
		log.Println(err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("false verificode"))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessExpiration := time.Duration(configs.Envs.JWTAccessExperationInMinites) * time.Minute
	accessToken, err := auth.CreateJWT(configs.Envs.JWTAccessSecret, u.Email, string(u.Role), accessExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshExpiration := time.Duration(configs.Envs.JWTRefreshExperationInHours) * time.Hour
	refreshToken, err := auth.CreateJWT(configs.Envs.JWTRefreshSecret, u.Email, string(u.Role), refreshExpiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Rdb.Set(context.Background(), u.Email, refreshToken, refreshExpiration).Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

func (s *Server) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	var payload ChangePasswordPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := s.Querier.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		log.Println(err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s didn't exists", payload.Email))
		return
	}

	val, err := s.Rdb.Get(context.Background(), payload.Email).Result()
	if err != nil || val != payload.VerifiCode {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("false verificode"))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = s.Querier.UpdatePasswordByEmail(r.Context(), store.UpdatePasswordByEmailParams{
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func GetIp(ip string) string {
	i := 0
	for j := len(ip) - 1; j >= 0; j-- {
		if ip[j] == ':' {
			i = j
			break
		}
	}

	if i == 0 {
		return "bad"
	}

	return ip[0:i]
}

func (s *Server) HandleSendVerifiCode(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	if len(ip) == 0 || len(ip) > 30 {
		utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("too many request"))
		return
	}

	ip = GetIp(ip)

	if err := s.Rdb.Get(context.Background(), ip).Err(); err == nil {
		utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("too many reqeust"))
		return
	}

	var payload SendVerifiCodePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	code := strconv.Itoa(rand.Intn(9000) + 1000) // 1000 ~ 9999

	if err := s.SendCode(payload.Email, code); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Rdb.Set(context.Background(), payload.Email, code, 2*time.Minute).Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Rdb.Set(context.Background(), ip, "wait for 1 mimnute", 1*time.Minute).Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "code had send",
	})
}

func (s *Server) SendCode(email, code string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", configs.Envs.MailUsername)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "验证码")

	html := fmt.Sprintf(`<div
        style={{
          display: "flex",
          flexDirection: "column",
          height: "174px",
          width: "174px",
          alignItems: "center",
          justifyContent: "center",
          borderColor: "black",
          marginLeft: "auto",
          marginRight: "auto",
          boxShadow: "inherit",
          backgroundColor: "#60a5fa",
          borderRadius: "7px",
        }}
      >
        <h1 style={{ textAlign: "center", color: "black" }}>您的验证码为:</h1>
        <p
          style={{
            textAlign: "center",
            color: "red",
            fontSize: "30px",
            lineHeight: "50px",
            fontWeight: 700,
          }}
        >
          %s
        </p>
      </div>`, code)
	m.SetBody("text/html", html)

	err := s.Mail.DialAndSend(m)

	return err
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
	UserName   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=3,max=130"`
	VerifiCode string `json:"verifi_code" validate:"required,min=4,max=4"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type GetUserResponse struct {
	Id        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

type ChangePasswordPayload struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
	VerifiCode  string `json:"verifi_code" validate:"required,min=4,max=4"`
}

type SendVerifiCodePayload struct {
	Email string `json:"email" validate:"required,email"`
}
