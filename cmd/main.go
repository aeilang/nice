package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aeilang/nice/auth"
	"github.com/aeilang/nice/configs"
	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

func main() {
	conf := configs.Envs
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.PublicHost, conf.Port, conf.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}

	mail := gomail.NewDialer(conf.MailHost, conf.MailPort, conf.MailUsername, conf.MailPassword)

	querier := store.New(db)
	serv := server.New(querier, rdb, mail)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(CORS)

	route(r, serv)

	http.ListenAndServe(":8888", r)
}

func route(r *chi.Mux, serv *server.Server) *chi.Mux {
	// public
	r.Group(func(r chi.Router) {
		r.Post("/login", serv.HandleLogin)
		r.Post("/register", serv.HandleRegister)
		r.Post("/verify", serv.HandleSendVerifiCode)
		r.Post("/forget", serv.HandleChangePassword)
		r.Post("/refresh", serv.HandleRefreshToken)
	})

	// protected
	r.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(serv.Querier))

		// User
		r.Route("/user", func(r chi.Router) {
			r.Get("/{id}", serv.HandleGetUser)
		})

	})

	return r
}

func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
