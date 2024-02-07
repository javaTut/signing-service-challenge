package repositories

import (
	"signing-service-challenge/crypto"
	"signing-service-challenge/domain"
	"signing-service-challenge/persistence"
	"testing"

	"github.com/google/uuid"
)

// var mutex sync.Mutex

func TestSaveDevice(t *testing.T) {
	var repository = createDeviceRepository()
	numOfDevices := 50
	createDevices(numOfDevices, repository)
	if repository.Count() != 50 {
		t.Errorf("%d devices saved, but %d expected", repository.Count(), 50)
	}
}

func TestGetDeviceById(t *testing.T) {
	var repository = createDeviceRepository()
	numOfDevices := 50
	ids := createDevices(numOfDevices, repository)
	for _, id := range ids {
		result, _ := repository.GetById(id)
		if result == nil {
			t.Error("device should be saved, but it isn't")
		}
		if result.Id != id {
			t.Errorf("got id %s, but %s expected", result.Id, id)
		}
	}
}

func TestGetDeviceByIdNotFound(t *testing.T) {
	var repository = createDeviceRepository()
	id := uuid.NewString()
	device := domain.NewSignatureDeviceWithoutKeys(id, crypto.RSA, "Test DEvice")
	repository.Save(*device)
	result, err := repository.GetById(id)
	if result == nil {
		t.Error("device should be found, but it is not")
	}
	resultNotFound, err := repository.GetById(id + " changed")
	if err == nil {
		t.Error("error should be returned, but it isn't")
	}
	if resultNotFound != nil {
		t.Error("device should not be found, but it is")
	}

}

func TestGetAllDevices(t *testing.T) {
	var repository = createDeviceRepository()
	numOfDevices := 50
	createDevices(numOfDevices, repository)
	if repository.Count() != numOfDevices {
		t.Errorf("got %d devices, %d expected", repository.Count(), numOfDevices)
	}
}

func TestDeleteDeviceById(t *testing.T) {
	var repository = createDeviceRepository()
	numOfDevices := 50
	ids := createDevices(numOfDevices, repository)
	if repository.Count() != numOfDevices {
		t.Errorf("%d devices should be saved, but there are %d", repository.Count(), numOfDevices)
	}
	for _, id := range ids {
		repository.DeleteById(id)
	}

	if repository.Count() != 0 {
		t.Errorf("%d devices left, but should be %d", repository.Count(), 0)
	}
}

func TestDeleteAllDevices(t *testing.T) {
	var repository = createDeviceRepository()
	numOfDevices := 50
	createDevices(numOfDevices, repository)
	if repository.Count() != numOfDevices {
		t.Errorf("%d devices should be saved, but there are %d", repository.Count(), numOfDevices)
	}
	repository.DeleteAll()

	if repository.Count() != 0 {
		t.Errorf("%d devices left, but should be %d", repository.Count(), 0)
	}
}

func createDeviceRepository() *SignatureDeviceInMemoryRepository {
	var db = persistence.NewInMemoryDB()
	return NewSignatureDeviceInMemoryRepository(db)
}

func createDevices(numOfDevices int, repository *SignatureDeviceInMemoryRepository) []string {
	ids := []string{}
	for i := 0; i < numOfDevices; i++ {
		id := uuid.NewString()
		ids = append(ids, id)
		device := domain.NewSignatureDeviceWithoutKeys(id, crypto.RSA, "Test DEvice")
		repository.Save(*device)
	}
	return ids
}
