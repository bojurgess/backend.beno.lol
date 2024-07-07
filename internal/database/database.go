package database

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Connection *sqlx.DB
}

type Meta struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID          string `db:"id"`
	DisplayName string `db:"display_name"`
	Email       string `db:"email"`
	Tokens
	Meta
}

type Tokens struct {
	AccessToken  string    `db:"access_token"`
	TokenType    string    `db:"token_type"`
	ExpiresAt    time.Time `db:"expires_at"`
	Scope        string    `db:"scope"`
	RefreshToken string    `db:"refresh_token"`
}

func (d *Database) Connect(addr string) {
	db, err := sqlx.Connect("sqlite3", addr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.MustExec(`PRAGMA foreign_keys = ON`)
	db.MustExec(`PRAGMA journal_mode = WAL`)
	db.MustExec(`PRAGMA foreign_keys = ON`)
	db.MustExec(`PRAGMA cache = shared`)
	db.MustExec(`PRAGMA mode = rwc`)

	d.Connection = db
}

func (d *Database) AddUser(u User) error {
	nstmt, err := d.Connection.PrepareNamed("INSERT INTO users (id, display_name, email, access_token, token_type, expires_at, scope, refresh_token) VALUES (:id, :display_name, :email, :access_token, :token_type, :expires_at, :scope, :refresh_token)")
	if err != nil {
		return err
	}

	_, err = nstmt.Exec(u)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateUser(u User) error {
	stmt := `UPDATE users SET display_name = COALESCE(:display_name, display_name), email = COALESCE(:email, email), access_token = COALESCE(:access_token, access_token), token_type = COALESCE(:token_type, token_type), expires_at = COALESCE(:expires_at, expires_at), scope = COALESCE(:scope, scope), refresh_token = COALESCE(:refresh_token, refresh_token) WHERE id = :id`
	_, err := d.Connection.NamedExec(stmt, u)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUser(id string) (User, error) {
	var u User
	stmt := `SELECT * FROM users WHERE id = ?`
	if err := d.Connection.Get(&u, stmt, id); err != nil {
		return u, err
	}

	return u, nil
}
