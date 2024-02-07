package repositories

import (
	"crypto/rand"
	"encoding/base64"
	"signing-service-challenge/domain"
	"signing-service-challenge/persistence"
	"testing"

	"github.com/google/uuid"
)

func TestSaveSignature(t *testing.T) {
	var repository = createSignatureRepository()
	numOfSignatures := 50
	createSignatures(numOfSignatures, repository)
	if repository.Count() != 50 {
		t.Errorf("%d signatures saved, but %d expected", repository.Count(), 50)
	}
}

func TestGetSignatureById(t *testing.T) {
	var repository = createSignatureRepository()
	numOfSignatures := 50
	ids := createSignatures(numOfSignatures, repository)
	for _, id := range ids {
		result, _ := repository.GetById(id)
		if result == nil {
			t.Error("signature should be saved, but it isn't")
		}
		if result.Id != id {
			t.Errorf("got id %s, expected %s", result.Id, id)
		}
	}
}

func TestGetSignatureByIdNotFound(t *testing.T) {
	var repository = createSignatureRepository()
	id := uuid.NewString()
	sig := generateRandomString(20)
	data := generateRandomString(20)
	deviceId := uuid.NewString()
	signature := domain.NewSignature(id, sig, data, deviceId)
	repository.Save(*signature)
	result, err := repository.GetById(id + " changed")
	if err == nil {
		t.Error("error should be returned, but it isn't")
	}
	if result != nil {
		t.Error("signature should not be found, but it is")
	}

}

func TestGetAllSignatures(t *testing.T) {
	var repository = createSignatureRepository()
	numOfSignatures := 50
	createSignatures(numOfSignatures, repository)
	signatures, _ := repository.GetAll()
	if len(signatures) != numOfSignatures {
		t.Errorf("got %d signatures, %d expected", len(signatures), numOfSignatures)
	}
}

func TestDeleteSignatureById(t *testing.T) {
	var repository = createSignatureRepository()
	numOfSignatures := 50
	ids := createSignatures(numOfSignatures, repository)

	if repository.Count() != numOfSignatures {
		t.Errorf("%d signatures should be saved, but there are %d", repository.Count(), numOfSignatures)
	}
	for _, id := range ids {
		repository.DeleteById(id)
	}
	if repository.Count() != 0 {
		t.Errorf("%d signatures left, but should be %d", repository.Count(), 0)
	}
}

func TestDeleteAllSignatures(t *testing.T) {
	var repository = createSignatureRepository()
	numOfSignatures := 50
	createSignatures(numOfSignatures, repository)

	if repository.Count() != numOfSignatures {
		t.Errorf("%d signatures should be saved, but there are %d", repository.Count(), numOfSignatures)
	}
	repository.DeleteAll()

	if repository.Count() != 0 {
		t.Errorf("%d signatures left, but should be %d", repository.Count(), 0)
	}
}

func createSignatureRepository() *SignatureInMemoryRepository {
	var db = persistence.NewInMemoryDB()
	return NewSignatureInMemoryRepository(db)
}

func createSignatures(numOfSignatures int, repository *SignatureInMemoryRepository) []string {
	ids := []string{}
	for i := 0; i < numOfSignatures; i++ {
		id := uuid.NewString()
		ids = append(ids, id)
		sig := generateRandomString(20)
		data := generateRandomString(20)
		deviceId := "ID12345"
		signature := domain.NewSignature(id, sig, data, deviceId)
		repository.Save(*signature)
	}
	return ids
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
