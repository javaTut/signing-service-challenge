package domain

import (
	"encoding/base64"
)

type SignatureDevice struct {
	Id               string
	Algorithm        string
	PrivateKey       []byte
	PublicKey        []byte
	Label            string
	SignatureCounter int
	LastSignature    string
}

func NewSignatureDeviceWithoutKeys(id string, algorithm string, label string) *SignatureDevice {
	return &SignatureDevice{
		Id:               id,
		Algorithm:        algorithm,
		Label:            label,
		SignatureCounter: 0,
		LastSignature:    base64.StdEncoding.EncodeToString([]byte(id)),
	}
}
