package services

import "errors"

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInactiveAccount        = errors.New("admin account is inactive")
	ErrNotFound               = errors.New("record not found")
	ErrInsufficientStock      = errors.New("insufficient ticket stock")
	ErrPaymentAlreadyVerified = errors.New("payment already verified")
	ErrInvalidSignature       = errors.New("invalid payment signature")
	ErrPaymentGatewayDisabled = errors.New("payment gateway is not configured")
	ErrOnlySuperAdmin         = errors.New("cannot remove the only active super admin")
)
