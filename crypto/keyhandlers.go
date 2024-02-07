package crypto

import (
	"errors"
	"signing-service-challenge/domain"
)

// supported algorithms
const (
	RSA = "RSA"
	ECC = "ECC"
)

// KeyPairHandler factory creates new instance based on the algorithm.
func GenerateKeyPairHandler(algorithm string) (KeyPairHandler, error) {
	switch algorithm {
	case RSA:
		return NewRSAKeyPairHandler(
			RSAGenerator{},
			NewRSAMarshaler(),
		), nil
	case ECC:
		return NewECCKeyPairHandler(
			ECCGenerator{},
			NewECCMarshaler(),
		), nil
	}
	return nil, errors.New("algorithm not supported")
}

// Wrappers around Generator and Marshaler
type KeyPairHandler interface {
	GenerateKeyPair() ([]byte, []byte, error)
	AttachKeyPair(*domain.SignatureDevice, []byte, []byte) (*domain.SignatureDevice, error)
	Unmarshal([]byte) (interface{}, error)
}

type RSAKeyPairHandler struct {
	Generator RSAGenerator
	Marshaler RSAMarshaler
}

func NewRSAKeyPairHandler(generator RSAGenerator, marshaler RSAMarshaler) *RSAKeyPairHandler {
	return &RSAKeyPairHandler{
		Generator: generator,
		Marshaler: marshaler,
	}
}

func (kph *RSAKeyPairHandler) GenerateKeyPair() ([]byte, []byte, error) {
	keyPair, err := kph.Generator.Generate()
	if err != nil {
		return nil, nil, err
	}
	marshaledPublicKey, marshaledPrivateKey, err := kph.Marshaler.Marshal(*keyPair)
	if err != nil {
		return nil, nil, err
	}
	return marshaledPrivateKey, marshaledPublicKey, nil
}

func (kph *RSAKeyPairHandler) AttachKeyPair(device *domain.SignatureDevice, privateKey []byte, publicKey []byte) (*domain.SignatureDevice, error) {
	if len(device.PrivateKey) > 0 || len(device.PublicKey) > 0 {
		return nil, domain.ErrKeyPairAlreadyAttached
	}
	device.PrivateKey = privateKey
	device.PublicKey = publicKey
	return device, nil
}

func (kph *RSAKeyPairHandler) Unmarshal(privateKey []byte) (interface{}, error) {
	key, err := kph.Marshaler.Unmarshal(privateKey)
	if err != nil {
		return nil, err
	}
	return key.Private, nil
}

type ECCKeyPairHandler struct {
	Generator ECCGenerator
	Marshaler ECCMarshaler
}

func NewECCKeyPairHandler(generator ECCGenerator, marshaler ECCMarshaler) *ECCKeyPairHandler {
	return &ECCKeyPairHandler{
		Generator: generator,
		Marshaler: marshaler,
	}
}

func (kph *ECCKeyPairHandler) GenerateKeyPair() ([]byte, []byte, error) {
	keyPair, err := kph.Generator.Generate()
	if err != nil {
		return nil, nil, err
	}
	marshaledPublicKey, marshaledPrivateKey, err := kph.Marshaler.Encode(*keyPair)
	if err != nil {
		return nil, nil, err
	}
	return marshaledPrivateKey, marshaledPublicKey, nil
}

func (kph *ECCKeyPairHandler) AttachKeyPair(device *domain.SignatureDevice, privateKey []byte, publicKey []byte) (*domain.SignatureDevice, error) {
	if len(device.PrivateKey) > 0 || len(device.PublicKey) > 0 {
		return nil, domain.ErrKeyPairAlreadyAttached
	}
	device.PrivateKey = privateKey
	device.PublicKey = publicKey
	return device, nil
}

func (kph *ECCKeyPairHandler) Unmarshal(privateKey []byte) (interface{}, error) {
	key, err := kph.Marshaler.Decode(privateKey)
	if err != nil {
		return nil, err
	}
	return key.Private, nil
}
