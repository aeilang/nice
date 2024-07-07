package server

import (
	"github.com/aeilang/nice/db/store"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type Server struct {
	Querier store.Querier
	Rdb     *redis.Client
	Mail    *gomail.Dialer
}

func New(queries *store.Queries, rdb *redis.Client, mail *gomail.Dialer) *Server {

	return &Server{
		Querier: queries,
		Rdb:     rdb,
		Mail:    mail,
	}
}
