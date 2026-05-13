package domain

import "errors"

// Sentinel errors for domain validation and repository operations.
var (
	ErrNotFound          = errors.New("caucion not found")
	ErrInvalidState      = errors.New("invalid estado value")
	ErrInvalidTransition = errors.New("invalid state transition")
)
