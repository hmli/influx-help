package influx_help

import (
	"go.uber.org/zap"
	"github.com/influxdata/influxdb/client/v2"
)

type DB struct {
	Addr string
	Username string
	Password string
	ShowSQL bool
	Logger *zap.Logger
	client client.Client
}

func NewDB(address, username, password string) *DB {
	logger, _ := zap.NewProduction() // TODO logger
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	db := DB {
		Addr: address,
		Username: username,
		Password: password,
		ShowSQL: false,
		Logger: logger,
		client: c,
	}
	return &db
}

func (db *DB) NewSession(database, precision string) (sess *Session) {
	return &Session{
		Database: database,
		Precision: precision,
		DB: db,
	}
}

type Session struct {
	DB *DB
	Database string
	Precision string
}

func (sess *Session) Measurement(m string) (stmt *Statement) {
	stmt = new(Statement)
	stmt.Init(sess)
	stmt.Measurement(m)
	return
}

func (sess *Session) Table(t string) (stmt *Statement) {
	return sess.Measurement(t)
}

