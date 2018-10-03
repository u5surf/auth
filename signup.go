// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

const (
	minPasswordLength = 8
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	// misc profile information
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	CompanyURL string `json:"companyUrl,omitempty"`
}

func addSignupRoutes(router *mux.Router, logger log.Logger, auth authable, userService userRepository) {
	router.Methods("POST").Path("/users/create").HandlerFunc(signupRoute(auth, userService))
}

func signupRoute(auth authable, userService userRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bs, err := read(r.Body)
		if err != nil {
			internalError(w, err, "signup")
			return
		}

		// read request body
		var signup signupRequest
		if err := json.Unmarshal(bs, &signup); err != nil {
			logger.Log("login", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Basic data sanity checks
		if err := validateEmail(signup.Email); err != nil {
			encodeError(w, err)
			return
		}
		if err := validatePassword(signup.Password); err != nil {
			encodeError(w, err)
			return
		}
		if err := validatePhone(signup.Phone); err != nil {
			encodeError(w, err)
			return
		}

		// find user
		u, err := userService.lookupByEmail(signup.Email)
		if err != nil && !strings.Contains(err.Error(), "user not found") {
			encodeError(w, errors.New("if this user exists, please try again with proper credentials"))
			return
		}
		if u == nil {
			var signup signupRequest
			if err := json.Unmarshal(bs, &signup); err != nil {
				encodeError(w, err)
				logger.Log("signup", fmt.Sprintf("failed parsing request json: %v", err))
				return
			}

			// store user
			userId := generateID()
			if userId == "" {
				internalError(w, fmt.Errorf("problem creating userId, err=%v", err), "signup")
				return
			}
			u = &User{
				ID:         userId,
				Email:      signup.Email,
				FirstName:  signup.FirstName,
				LastName:   signup.LastName,
				Phone:      signup.Phone,
				CompanyURL: signup.CompanyURL,
				CreatedAt:  time.Now(),
			}
			if err := userService.upsert(u); err != nil {
				internalError(w, fmt.Errorf("problem writing user: %v", err), "signup")
				return
			}

			if err := auth.writePassword(u.ID, signup.Password); err != nil {
				internalError(w, fmt.Errorf("problem writing user credentials: %v", err), "signup")
				return
			}

			// signup worked, yay!
			// TODO(adam): email approval link and clickthrough
		} else {
			// user found, so reject signup
			encodeError(w, errors.New("user already exists"))
		}
	}
}

func validateEmail(email string) error {
	if email == "" || !strings.Contains(email, "@") {
		return errors.New("no email provided")
	}
	return nil
}

func validatePassword(pass string) error {
	if pass == "" {
		return errors.New("no password provided")
	}
	if n := utf8.RuneCountInString(pass); n < minPasswordLength {
		return fmt.Errorf("password required to be at least %d characters", n)
	}
	return nil
}

func validatePhone(phone string) error {
	phone = strings.Replace(phone, "-", "", -1)
	phone = strings.Replace(phone, ".", "", -1)
	phone = strings.Replace(phone, " ", "", -1)
	if m, _ := regexp.MatchString("^\\+?[1-9]\\d{1,14}$", phone); !m {
		return errors.New("phone number is invalid")
	}
	return nil
}
