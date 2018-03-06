package influx_help

import "go.uber.org/zap"

type DB struct {
	Addr string
	Username string
	Password string
	ShowSQL bool
	Logger *zap.Logger
}

func NewDB(address, username, password string) *DB {
	logger, _ := zap.NewProduction()
	db := DB {
		Addr: address,
		Username: username,
		Password: password,
		ShowSQL: false,
		Logger: logger,
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
	stmt.Session = sess
	stmt.Init(sess)
	return
}

func (sess *Session) Table(t string) (stmt *Statement) {
	return sess.Measurement(t)
}

