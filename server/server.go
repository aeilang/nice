package server

import (
	"github.com/aeilang/nice/db/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Querier store.Querier
	Pool    *pgxpool.Pool
}

func New(queries *store.Queries, pool *pgxpool.Pool) *Server {
	return &Server{
		Querier: queries,
		Pool:    pool,
	}
}
