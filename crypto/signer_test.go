package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"reflect"
	"testing"
)

func TestSignerFactoryShouldReturnRSASigner(t *testing.T) {
	pk := rsa.PrivateKey{}
	signer, _ := GenerateSigner(&pk)
	if reflect.TypeOf(signer) != reflect.TypeOf(&RSASigner{}) {
		t.Error("expected RSASigner, but another returned")
	}
}

func TestSignerFactoryShouldReturnECCSigner(t *testing.T) {
	pk := ecdsa.PrivateKey{}
	signer, _ := GenerateSigner(&pk)
	if reflect.TypeOf(signer) != reflect.TypeOf(&ECCSigner{}) {
		t.Error("expected ECCSigner, but another returned")
	}
}

func TestSignerFactoryShouldReturnNilAndErrorForNotSupportedPrivateKeys(t *testing.T) {
	pk := struct{}{} // unsupported key
	signer, err := GenerateSigner(&pk)
	if signer != nil {
		t.Error("it should return nil for unsupported private keys")
	}
	if err == nil {
		t.Error("it should return error for unsupported private keys")
	}

}

func TestRSASignerSameMessageShouldPass(t *testing.T) {
	generator := RSAGenerator{}
	keys, _ := generator.Generate()
	signer := NewRSASigner(keys.Private)
	data := []byte("Data for signing")
	signature, _ := signer.Sign(data)
	verified, _ := signer.Verify(data, signature)
	if !verified {
		t.Errorf("got verified: %t, but should be true", verified)
	}
}

func TestRSASignerTemperedMessageShouldFail(t *testing.T) {
	generator := RSAGenerator{}
	keys, _ := generator.Generate()
	signer := NewRSASigner(keys.Private)
	data := []byte("Data for signing")
	dataTemp := []byte("Data for signing (tempered)")
	signature, _ := signer.Sign(data)
	verified, _ := signer.Verify(dataTemp, signature)
	if verified {
		t.Errorf("got verified: %t, but should be true", verified)
	}
}

func TestECCSignerSameMessageShouldPass(t *testing.T) {
	generator := ECCGenerator{}
	keys, _ := generator.Generate()
	signer := NewECCSigner(keys.Private)
	data := []byte("Data for signing")
	signature, _ := signer.Sign(data)
	verified, _ := signer.Verify(data, signature)
	if !verified {
		t.Errorf("got verified: %t, but should be true", verified)
	}
}

func TestECCSignerTemperedMessageShouldFail(t *testing.T) {
	generator := ECCGenerator{}
	keys, _ := generator.Generate()
	signer := NewECCSigner(keys.Private)
	data := []byte("Data for signing")
	dataTemp := []byte("Data for signing (tempered)")
	signature, _ := signer.Sign(data)
	verified, _ := signer.Verify(dataTemp, signature)
	if verified {
		t.Errorf("got verified: %t, but should be true", verified)
	}
}
