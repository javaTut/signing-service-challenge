package domain

import "errors"

var (
	ErrDeviceNotFound         = errors.New("device not found")
	ErrSignatureNotFound      = errors.New("signature not found")
	ErrKeyPairAlreadyAttached = errors.New("key pair for this device already attached")
)
