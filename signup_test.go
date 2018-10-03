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

func TestSignup__phone(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"", false},
		{"abcdefgh", false},
		{"a0a0a0a0", false},
		{"1a1a1a1a", false},
		{"0123456789", false},
		{"1", false},
		{"10", true},
		{"109", true},
		{"10909", true},
		{"109090", true},
		{"2090999", true},
		{"30909999", true},
		{"409099999", true},
		{"5090999999", true},
		{"60909999999", true},
		{"709099999999", true},
		{"8090999999999", true},
		{"90909999999999", true},
		{"1009099999999999", false},
		{"+14155552671000", true},
		{"+1415555267100001", false},
		{"+1-6174443000", true},
		{"+33 1 5669 6201", true},
		{"49-8994006308", true},
		{"+972-732858700", true},
		{"+81-90-1234-5678", true},
		{"666.666.6666", true},
	}
	for i := range cases {
		err := validatePhone(cases[i].input)
		if cases[i].valid && err == nil {
			continue // valid
		}
		if !cases[i].valid && err != nil {
			continue // known bad
		}
		t.Errorf("input=%q, err=%v", cases[i].input, err)
	}
}
