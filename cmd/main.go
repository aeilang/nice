package main

import (
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
)

func main() {
	conf := configs.Envs
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.PublicHost, conf.Port, conf.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	querier := store.New(db)
	serv := server.New(querier)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	route(r, serv)

	http.ListenAndServe(":3000", r)
}

func route(r *chi.Mux, serv *server.Server) *chi.Mux {
	// public
	r.Group(func(r chi.Router) {
		r.Post("/login", serv.HandleLogin)
		r.Post("/regis", serv.HandleRegister)
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
