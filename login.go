// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func addLoginRoutes(router *mux.Router, logger log.Logger, auth authable, userService userRepository) {
	router.Methods("GET").Path("/users/login").HandlerFunc(checkLogin(logger, auth, userService))
	router.Methods("POST").Path("/users/login").HandlerFunc(loginRoute(logger, auth, userService))
}

func checkLogin(logger log.Logger, auth authable, userService userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		cookie := extractCookie(r)
		if cookie == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		userId, err := auth.findUserId(cookie.Value)
		if err != nil {
			internalError(w, err, "login")
			return
		}
		if userId == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("X-User-Id", userId)
		w.WriteHeader(http.StatusOK)
	}
}

func loginRoute(logger log.Logger, auth authable, userService userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bs, err := read(r.Body)
		if err != nil {
			internalError(w, err, "login")
			return
		}

		// read request body
		var login loginRequest
		if err := json.Unmarshal(bs, &login); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Basic data sanity checks
		if err := validateEmail(login.Email); err != nil {
			encodeError(w, err)
			return
		}
		if err := validatePassword(login.Password); err != nil {
			encodeError(w, err)
			return
		}

		// find user by email
		u, err := userService.lookupByEmail(login.Email)
		if err != nil {
			// Mark this (and password check) as failure only because
			// the user is involved at this point. Otherwise it's their
			// developer's problem (i.e. bad json).
			authFailures.With("method", "web").Add(1)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// find user by userId and password
		if err := auth.checkPassword(u.ID, login.Password); err != nil {
			authFailures.With("method", "web").Add(1)
			logger.Log("login", fmt.Sprintf("userId=%s failed: %v", u.ID, err))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// success route, let's finish!
		authSuccesses.With("method", "web").Add(1)
		cookie, err := createCookie(u.ID, auth)
		if err != nil {
			internalError(w, err, "login")
			return
		}
		if cookie == nil {
			logger.Log("login", fmt.Sprintf("nil cookie for userId=%s", u.ID))
			internalError(w, err, "login")
			return
		}
		if err := auth.writeCookie(u.ID, cookie); err != nil {
			internalError(w, err, "login")
			return
		}

		http.SetCookie(w, cookie)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			internalError(w, err, "login")
			return
		}
	}
}
