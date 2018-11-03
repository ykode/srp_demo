package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/ykode/srp_demo/server/internal/domain"
	"github.com/ykode/srp_demo/server/internal/query"
)

type SSLMode string

const (
	SSLDisable    SSLMode = "disable"
	SSLRequire    SSLMode = "require"
	SSLVerifyCA   SSLMode = "verify-ca"
	SSLVerifyFull SSLMode = "verify-full"
)

type PostgreSQLStorage struct {
	db  *sql.DB
	log *logrus.Entry
}

func NewPostgreSQLStorage(host, username, password, dbname string, port int, sslmode SSLMode) (*PostgreSQLStorage, error) {
	log := logrus.WithFields(logrus.Fields{
		"topic": "postgresql",
	})

	connStr := fmt.Sprintf("user=%s password=%s, dbname=%s host=%s port=%d sslmode=%s",
		username, password, dbname, host, port, sslmode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &PostgreSQLStorage{db: db, log: log}, nil
}

func (s *PostgreSQLStorage) _findIdentityByUserName(username string) (*domain.Identity, error) {
	q := `SELECT "username", "salt", "verifier" FROM identities WHERE username=$1`

	var rows *sql.Rows
	var err error

	if rows, err = s.db.Query(q, username); err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, query.ErrNotFound
	}

	var uname string
	var salt, verifier []byte

	if err = rows.Scan(&username, &salt, &verifier); err != nil {
		return nil, err
	}

	var id *domain.Identity

	if id, err = domain.NewIdentity(uname, salt, verifier); err != nil {
		return nil, err
	}

	return id, nil

}

func (s *PostgreSQLStorage) _saveIdentity(id *domain.Identity) error {
	q := `INSERT INTO identites(username, salt, verifier) VALUES ($1, $2, $3)
				ON CONFLICT(username) DO
				UPDATE identities SET username = $1, salt = $2, verifier = $3`

	if _, err := s.db.Exec(q); err != nil {
		return err
	}

	return nil
}

func (s *PostgreSQLStorage) _findSessionById(sessionId uuid.UUID) (*domain.Session, error) {
	q := `SELECT "id", "salt", "master_key", "state", "v" FROM sessions WHERE "id" = $1`

	var rows *sql.Rows
	var err error

	if rows, err = s.db.Query(q, sessionId); err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, query.ErrNotFound
	}

	var id uuid.UUID
	var salt, masterKey, v []byte
	var state domain.SessionState

	if err = rows.Scan(&id, &salt, &masterKey, &state); err != nil {
		return nil, err
	}

	var session *domain.Session

	if session, err = domain.BuildSession(id, salt, masterKey, v, state); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *PostgreSQLStorage) _saveSession(*domain.Session) error {
	q := `INSERT INTO sessions("id", "salt", "master_key", "state", "v")( $1, $2, $3, $4, $5)
				ON CONFLICT("id")
				UPDATE sessions SET "id" = $1, "salt" = $2, "master_key" = $3, "state" = $4, "v" = $5)`

	if _, err := s.db.Exec(q); err != nil {
		return err
	}

	return nil
}

func (s *PostgreSQLStorage) FindIdentityByUserName(username string) <-chan query.Result {
	c := make(chan query.Result)

	go func() {
		id, err := s._findIdentityByUserName(username)

		c <- query.Result{Result: id, Err: err}

	}()

	return c
}

func (s *PostgreSQLStorage) SaveIdentity(id *domain.Identity) <-chan error {
	c := make(chan error)

	go func() {
		c <- s._saveIdentity(id)
	}()

	return c
}

func (s *PostgreSQLStorage) FindSessionById(sessionId uuid.UUID) <-chan query.Result {
	c := make(chan query.Result)

	go func() {
		sess, err := s._findSessionById(sessionId)

		c <- query.Result{Result: sess, Err: err}
	}()

	return c
}

func (s *PostgreSQLStorage) SaveSession(sess *domain.Session) <-chan error {
	c := make(chan error)

	go func() {
		c <- s._saveSession(sess)
	}()

	return c
}
