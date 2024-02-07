package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	// "signing-service-code-challenge/services"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	Verify(signedData, signature []byte) (bool, error) // for testing purposes
}

// Signer factory generates new Signer instance based on the data type of the passed private key.
func GenerateSigner(s interface{}) (Signer, error) {
	switch x := s.(type) {
	case *rsa.PrivateKey:
		return NewRSASigner(x), nil
	case *ecdsa.PrivateKey:
		return NewECCSigner(x), nil
	}
	return nil, errors.New("algorithm for the given private key type not supported!")
}

type RSASigner struct {
	privateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		privateKey: privateKey,
	}
}

func (s RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(dataToBeSigned)
	if err != nil {
		return []byte{}, err
	}
	hashSum := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashSum)
	if err != nil {
		return []byte{}, err
	}
	return signature, nil
}

func (s RSASigner) Verify(signedData []byte, signature []byte) (bool, error) {
	hash := sha256.New()
	_, err := hash.Write(signedData)
	if err != nil {
		return false, err
	}
	hashSum := hash.Sum(nil)
	err = rsa.VerifyPKCS1v15(&s.privateKey.PublicKey, crypto.SHA256, hashSum, signature)
	if err != nil {
		return false, nil
	}
	return true, nil
}

type ECCSigner struct {
	privateKey *ecdsa.PrivateKey
}

func NewECCSigner(privateKey *ecdsa.PrivateKey) *ECCSigner {
	return &ECCSigner{
		privateKey: privateKey,
	}
}

func (s ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	signature, err := ecdsa.SignASN1(rand.Reader, s.privateKey, hash[:])

	if err != nil {
		return []byte{}, err
	}
	return signature, nil
}

func (s ECCSigner) Verify(signedData []byte, signature []byte) (bool, error) {
	hash := sha256.Sum256(signedData)
	return ecdsa.VerifyASN1(&s.privateKey.PublicKey, hash[:], signature), nil

}
