package models

import "sync"

var InvalidTokens = struct {
	sync.RWMutex
	Tokens map[string]bool
}{Tokens: make(map[string]bool)}
