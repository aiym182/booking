package dbrepo

import (
	"database/sql"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/repository"
)

type PostgresDbRepo struct {
	App *config.Config
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.Config
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.Config) repository.DatabaseRepo {

	return &PostgresDbRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.Config) repository.DatabaseRepo {

	return &testDBRepo{
		App: a,
	}
}
