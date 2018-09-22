package main

import (
	"testing"
)

func TestSignup__email(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"", false},
		{"test@moov.io", true},
	}
	for i := range cases {
		err := validateEmail(cases[i].input)
		if cases[i].valid && err == nil {
			continue // valid
		}
		if !cases[i].valid && err != nil {
			continue // known bad
		}
		t.Errorf("input=%q, err=%v", cases[i].input, err)
	}
}

func TestSignup__pass(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"", false},
		{"superlongpassword", true},
	}
	for i := range cases {
		err := validatePassword(cases[i].input)
		if cases[i].valid && err == nil {
			continue // valid
		}
		if !cases[i].valid && err != nil {
			continue // known bad
		}
		t.Errorf("input=%q, err=%v", cases[i].input, err)
	}
}
