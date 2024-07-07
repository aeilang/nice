package server

import (
	"github.com/aeilang/nice/db/store"
)

type Server struct {
	Querier store.Querier
}

func New(queries *store.Queries) *Server {
	return &Server{
		Querier: queries,
	}
}
