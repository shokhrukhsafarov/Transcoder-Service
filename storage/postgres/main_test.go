package postgres_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/pkg/db"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/storage"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/manveru/faker"
	"github.com/stretchr/testify/assert"
)

var (
	err          error
	cfg          config.Config
	strg         storage.StorageI
	fakeData     *faker.Faker
	postgresConn *sqlx.DB
	l            *logger.Logger
)

func CreateRandomId(t *testing.T) string {
	id, err := uuid.NewRandom()
	assert.NoError(t, err)
	return id.String()
}

func TestMain(m *testing.M) {
	cfg = config.Load()
	l = logger.New(cfg.LogLevel)

	conStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
		"disable",
	)
	fakeData, _ = faker.New("en")
	postgresConn, err = sqlx.Open("postgres", conStr)

	if err != nil {
		log.Fatal(err)
	}

	pg := &db.Postgres{}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	pg.Db = postgresConn

	strg = storage.New(pg, nil, cfg)

	fakeData, _ = faker.New("en")

	os.Exit(m.Run())
}
