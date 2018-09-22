// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
)

// TODO(adam): validate email (by sending with approval_code) [probably diff file]

// TODO(adam): having a row in user_approval_codes means user isn't verified. Delete on email approval.

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Phone      string    `json:"phone"`
	CompanyURL string    `json:"companyUrl"`
	CreatedAt  time.Time `json:"createdAt"`
}

var (
	dropPlusExtender = regexp.MustCompile(`(\+.*)$`)
	dropPeriods      = strings.NewReplacer(".", "")
)

const (
	bcryptCostFactor = 10

	// from 'go doc time Time.String', format used in sqlite columns
	serializedTimestampFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func (u *User) cleanEmail() string {
	return cleanEmail(u.Email)
}

// cleanEmail strips all the funky characters from an email address.
//
// Essentially this boils down to the following pattern (ignoring case):
//   [a-z0-9]@[a-z0-9].[a-z]
//
// Callers should be aware of when an empty string is returned.
func cleanEmail(email string) string {
	// split at '@'
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}

	parts[0] = dropPlusExtender.ReplaceAllString(parts[0], "")
	parts[0] = dropPeriods.Replace(parts[0])

	return strings.ToLower(strings.Join(parts, "@"))
}

// generateID creates a new ID for our auth system.
// Do no assume anything about these ID's other than
// they are strings. Case matters
func generateID() string {
	bs := make([]byte, 20)
	n, err := rand.Read(bs)
	if err != nil || n == 0 {
		logger.Log("generateID", fmt.Sprintf("n=%d, err=%v", n, err))
		return ""
	}
	return strings.ToLower(hex.EncodeToString(bs))
}

type userRepository interface {
	lookupByUserId(id string) (*User, error)

	// lookupByEmail finds a user by the given email address.
	// It's recommended you provide User.cleanEmail()
	lookupByEmail(email string) (*User, error)

	upsert(*User) error
}

type sqliteUserRepository struct {
	db  *sql.DB
	log log.Logger
}

func (s *sqliteUserRepository) lookupByUserId(userId string) (*User, error) {
	query := `select u.email, u.created_at, ud.first_name, ud.last_name, ud.phone, ud.company_url
from users as u
inner join user_details as ud
on u.user_id = ud.user_id
where u.user_id = ?
limit 1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(userId)

	u := &User{}
	u.ID = userId
	var createdAt string // needs parsing
	row.Scan(&u.Email, &createdAt, &u.FirstName, &u.LastName, &u.Phone, &u.CompanyURL)
	t, err := time.Parse(serializedTimestampFormat, createdAt)
	if err != nil {
		s.log.Log("user", fmt.Sprintf("bad users.created_at format %q: %v", createdAt, err))
	}
	u.CreatedAt = t
	return u, nil
}

func (s *sqliteUserRepository) lookupByEmail(email string) (*User, error) {
	query := `select user_id from users where clean_email = ?`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(cleanEmail(email))

	var userId string
	row.Scan(&userId)

	if userId == "" {
		return nil, errors.New("user not found")
	}
	return s.lookupByUserId(userId)
}

func (s *sqliteUserRepository) upsert(inc *User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// insert/update into 'users'
	query := `replace into users (user_id, email, clean_email, created_at) values (?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		e := tx.Rollback()
		return fmt.Errorf("problem preparing users query userId=%s, err=%v, rollback err=%v", inc.ID, err, e)
	}
	_, err = stmt.Exec(inc.ID, inc.Email, inc.cleanEmail(), inc.CreatedAt.Format(serializedTimestampFormat))
	if err != nil {
		e := tx.Rollback()
		return fmt.Errorf("problem upserting users userId=%s, err=%v, rollback err=%v", inc.ID, err, e)
	}

	// insert/update into 'user_details'
	query = `replace into user_details (user_id, first_name, last_name, phone, company_url) values (?, ?, ?, ?, ?)`
	stmt, err = tx.Prepare(query)
	if err != nil {
		e := tx.Rollback()
		return fmt.Errorf("problem preparing user_details query userId=%s, err=%v, rollback err=%v", inc.ID, err, e)
	}
	_, err = stmt.Exec(inc.ID, inc.FirstName, inc.LastName, inc.Phone, inc.CompanyURL)
	if err != nil {
		e := tx.Rollback()
		return fmt.Errorf("problem upserting user_details userId=%s, err=%v, rollback err=%v", inc.ID, err, e)
	}

	return tx.Commit()
}

// authable represents the interactions of a user's authentication
// status. This boils down to password comparison and cookie data.
type authable interface {
	findUserId(data string) (string, error)
	invalidateCookies(userId string) error
	writeCookie(userId string, cookie *http.Cookie) error

	// checkPassword compares the provided password for the user.
	// a non-nil error is returned if the passwords don't match
	// or that the userId doesn't exist.
	checkPassword(userId string, pass string) error
	writePassword(userId string, pass string) error
}

type auth struct {
	db  *sql.DB
	log log.Logger
}

// findUserId takes cookie data and returns the userId associated
func (a *auth) findUserId(data string) (string, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return "", errors.New("no cookie data provided")
	}

	// the SHA256 checksum is stored, not the actual data.
	data, err := hash(data)
	if err != nil {
		return "", err
	}

	query := `select user_id from user_cookies where data == ? and valid_until > ?`
	stmt, err := a.db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(data, time.Now().Format(serializedTimestampFormat))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var userId string
	for rows.Next() {
		rows.Scan(&userId)
		if userId != "" {
			return userId, nil
		}
	}
	return "", nil
}

func (a *auth) invalidateCookies(userId string) error {
	stmt, err := a.db.Prepare(`delete from user_cookies where user_id = ?`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(userId)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	a.log.Log("user", fmt.Sprintf("deleted %d cookies for userId=%s", n, userId))
	return nil
}

func (a *auth) writeCookie(userId string, cookie *http.Cookie) error {
	query := `insert or replace into user_cookies (user_id, data, valid_until) values (?, ?, ?)`
	stmt, err := a.db.Prepare(query)
	if err != nil {
		return err
	}

	// hash the data
	data, err := hash(cookie.Value)
	if err != nil {
		return err
	}
	validUntil := cookie.Expires.Format(serializedTimestampFormat)

	// write row
	stmt.Exec(userId, data, validUntil)
	return nil
}

// fakeBcryptRounds just performs a bcrypt.GenerateFromPassword and then
// a bcrypto.CompareHashAndPassword afterwords. In an attempt to make happy
// and sad paths take "approximately" the same.
func (a *auth) fakeBcryptRounds() {
	id := generateID()
	bcrypt.CompareHashAndPassword([]byte(id), []byte(id))
}

func (a *auth) checkPassword(userId string, incoming string) error {
	stmt, err := a.db.Prepare(`select password, salt from user_passwords where user_id = ?`)
	if err != nil {
		a.fakeBcryptRounds()
		return err
	}

	row := stmt.QueryRow(userId)
	var storedPassword, storedSalt string
	if err := row.Scan(&storedPassword, &storedSalt); err != nil {
		a.fakeBcryptRounds()
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(incoming+storedSalt))
}

// writePassword saves a user's password. This function performs no authn/z and generates a new
// salt (which is also saved).
func (a *auth) writePassword(userId string, pass string) error {
	salt := generateID()
	incomingPassword, err := bcrypt.GenerateFromPassword([]byte(pass+salt), bcryptCostFactor)
	if err != nil {
		return err
	}

	stmt, err := a.db.Prepare(`replace into user_passwords (user_id, password, salt) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId, string(incomingPassword), salt)
	if err != nil {
		return err
	}
	a.log.Log("user", fmt.Sprintf("userId=%s updated password", userId))
	return nil
}

func hash(in string) (string, error) {
	ss := sha256.New()
	n, err := ss.Write([]byte(in))
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", err
	}
	return hex.EncodeToString(ss.Sum(nil)), nil
}
