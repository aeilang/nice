package main

import (
	"context"
	"net/http"

	"github.com/aeilang/nice/auth"
	"github.com/aeilang/nice/configs"
	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbConfig, err := configs.DBConfig()
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		panic(err)
	}
	defer pool.Close()

	querier := store.New(pool)
	serv := server.New(querier, pool)

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

		// product
		r.Route("/products", func(r chi.Router) {
			r.Get("/{id}", serv.HandleGetProduct)
			r.Get("/", serv.HandleGetProducts)
			r.Post("/", serv.HandleCreateProduct)

		})

		// order
		r.Route("/order", func(r chi.Router) {
			r.Post("/", serv.HandleCheckout)
		})
	})

	return r
}
