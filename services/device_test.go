package services

import (
	"fmt"
	"signing-service-challenge/crypto"
	"signing-service-challenge/dto"
	"signing-service-challenge/lockers"
	"signing-service-challenge/persistence"
	"signing-service-challenge/repositories"
	"sync"
	"testing"

	"github.com/google/uuid"
	// "time"
)

var db = persistence.NewInMemoryDB()
var repository = repositories.NewSignatureDeviceInMemoryRepository(db)
var mutex sync.Mutex

var locker = lockers.NewDeviceLockerWithGlobalMapProtection(&mutex)

func TestCreateSignatureDevice(t *testing.T) {
	service := NewSignatureDeviceService(repository, locker)
	id := uuid.NewString()
	algorithm := crypto.RSA
	label := "First Device"
	device, _ := service.CreateSignatureDevice(id, algorithm, label)
	if device.Id != id {
		t.Errorf("got ID: %s expected: %s.", device.Id, id)
	}
	if device.Algorithm != algorithm {
		t.Errorf("got algorithm: %s expected: %s.", device.Algorithm, algorithm)
	}
	deviceFromDb, _ := repository.GetById(id)
	if deviceFromDb == nil {
		t.Error("device should be saved into the database, but it isn't.")
	}
	if deviceFromDb.PrivateKey == nil || deviceFromDb.PublicKey == nil {
		t.Error("key pair should be set, but it isn't.")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestCreateSignatureDeviceForUnsupportedAlgorithmShouldFail(t *testing.T) {
	service := NewSignatureDeviceService(repository, locker)
	id := uuid.NewString()
	algorithm := "UNSUPPORTED"
	label := "First Device"
	device, err := service.CreateSignatureDevice(id, algorithm, label)
	if device != nil {
		t.Error("device should not be created")
	}
	if err == nil {
		t.Error("should throw error")
	}
	deviceFromDb, _ := repository.GetById(id)
	if deviceFromDb != nil {
		t.Error("device should not be saved into db.")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestRSASigningDataWithTheSameDeviceConcurrently(t *testing.T) {
	// assumption that one device could be used from multiple clients concurrently
	testSigningDataOneDeviceMultipleClientsConcurrently(t, crypto.RSA, locker, 1000)
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestECCSigningDataWithTheSameDeviceConcurrently(t *testing.T) {
	// assumption that one device could be used from multiple clients concurrently
	testSigningDataOneDeviceMultipleClientsConcurrently(t, crypto.ECC, locker, 1000)
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestSigningDataByMultipleDevicesConcurrentlyEachDeviceUsedOnlyOnce(t *testing.T) {
	service := NewSignatureDeviceService(repository, locker)

	createDevices(100, crypto.RSA, *service)
	createDevices(100, crypto.ECC, *service)

	allDevices, err := service.GetAll()
	if err != nil {
		t.Fatal("no devices in db, test could not continue")
	}

	// divide devices in two separated slices
	mid := len(allDevices) / 2
	firstHalf := allDevices[:mid]
	secondHalf := allDevices[mid:]
	// execute sign transactions concurrently
	var wg sync.WaitGroup
	testSigningDataByMultipleDevicesOnlyOneSignatureConcurrently(t, &wg, firstHalf, *service)
	testSigningDataByMultipleDevicesOnlyOneSignatureConcurrently(t, &wg, secondHalf, *service)
	wg.Wait()
	t.Cleanup(func() {
		repository.DeleteAll()
	})

}

func TestSigningDataByMultipleDevicesConcurrentlyMultipleSignaturesPerDevice(t *testing.T) {
	service := NewSignatureDeviceService(repository, locker)

	createDevices(100, crypto.RSA, *service)
	createDevices(100, crypto.ECC, *service)

	allDevices, err := service.GetAll()
	if err != nil {
		t.Fatal("no devices in db, test could not continue")
	}

	// divide devices in two separated slices
	mid := len(allDevices) / 2
	firstHalf := allDevices[:mid]
	secondHalf := allDevices[mid:]
	messages := [5]string{
		"message 1",
		"message 2",
		"message 3",
		"message 4",
		"message 5",
	}
	// execute sign transactions concurrently
	var wg sync.WaitGroup
	testSigningDataByMultipleDevicesMultipleSignaturesConcurrently(t, &wg, firstHalf, *service, messages[:])
	testSigningDataByMultipleDevicesMultipleSignaturesConcurrently(t, &wg, secondHalf, *service, messages[:])
	wg.Wait()
	devicesAfterSigning, _ := service.GetAll()
	// each device should sign 5 messages
	for _, d := range devicesAfterSigning {
		if d.SignatureCounter != len(messages) {
			t.Errorf("device should sign %d messages, got %d instead", len(messages), d.SignatureCounter)
		}
	}

}

func TestRSASignatureVerificationShouldReturnTrue(t *testing.T) {
	verified := testSignatureVerification(t, crypto.RSA, "")
	if !verified {
		t.Error("signature should be verified, but it failed")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestRSASignatureVerificationShouldReturnFalse(t *testing.T) {
	verified := testSignatureVerification(t, crypto.RSA, "tempered")
	if verified {
		t.Error("verification should fail, but it succeeded")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestECCSignatureVerificationShouldReturnTrue(t *testing.T) {
	verified := testSignatureVerification(t, crypto.ECC, "")
	if !verified {
		t.Error("signature should be verified, but it failed")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func TestECCSignatureVerificationShouldReturnFalse(t *testing.T) {
	verified := testSignatureVerification(t, crypto.ECC, "tempered")
	if verified {
		t.Error("verification should fail, but it succeeded")
	}
	t.Cleanup(func() {
		repository.DeleteAll()
	})
}

func testSigningDataByMultipleDevicesOnlyOneSignatureConcurrently(
	t *testing.T,
	wg *sync.WaitGroup,
	devices []dto.SignatureDeviceResponse,
	service SignatureDeviceService,
) {
	data := "message to be signed"
	for _, device := range devices {
		wg.Add(1)
		go func(d dto.SignatureDeviceResponse) {
			defer wg.Done()
			deviceBeforeSigning, _ := service.GetById(d.Id)
			result, _ := service.SignTransaction(d.Id, data)
			deviceAfterSigning, _ := service.GetById(d.Id)
			if deviceBeforeSigning.LastSignature == deviceAfterSigning.LastSignature {
				t.Error("last signature value should be different after each sign operation")
			}
			if (deviceAfterSigning.SignatureCounter - deviceBeforeSigning.SignatureCounter) != 1 {
				t.Errorf(
					"signature counter should increase by 1, increased by %d instead",
					(deviceAfterSigning.SignatureCounter - deviceBeforeSigning.SignatureCounter),
				)
			}
			if deviceAfterSigning.SignatureCounter != 1 {
				t.Errorf(
					"only one message should be signed, but got %d instead",
					deviceAfterSigning.SignatureCounter,
				)
			}
			if result == nil {
				t.Error("signature should be created, but it isn't")
			}
		}(device)
	}
}

func testSigningDataByMultipleDevicesMultipleSignaturesConcurrently(
	t *testing.T,
	wg *sync.WaitGroup,
	devices []dto.SignatureDeviceResponse,
	service SignatureDeviceService,
	data []string,
) {

	for _, device := range devices {
		wg.Add(1)
		go func(d dto.SignatureDeviceResponse) {
			defer wg.Done()
			for _, m := range data {
				wg.Add(1)
				go func(m string) {
					defer wg.Done()
					deviceBeforeSigning, _ := service.GetById(d.Id)
					result, _ := service.SignTransaction(d.Id, m)
					deviceAfterSigning, _ := service.GetById(d.Id)
					if deviceBeforeSigning.LastSignature == deviceAfterSigning.LastSignature {
						t.Error("last signature value should be different after each sign operation")
					}
					if result == nil {
						t.Error("signature should be created, but it isn't")
					}
				}(m)
			}
		}(device)
	}
}

func createDevices(n int, algorithm string, service SignatureDeviceService) {
	for i := 0; i < n; i++ {
		id := uuid.NewString()
		label := fmt.Sprintf("%s Device %d", algorithm, i)
		service.CreateSignatureDevice(id, algorithm, label)
	}
}

func testSigningDataOneDeviceMultipleClientsConcurrently(t *testing.T, algorithm string, locker lockers.DeviceLocker, numOfSignatures int) {
	service := NewSignatureDeviceService(repository, locker)
	id := uuid.NewString()
	label := "First Device"
	data := "message to be signed"
	service.CreateSignatureDevice(id, algorithm, label)
	var wg sync.WaitGroup
	// execute signing concurrently
	for i := 0; i < numOfSignatures; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deviceBeforeSigning, _ := service.GetById(id)
			result, _ := service.SignTransaction(id, data)
			deviceAfterSigning, _ := service.GetById(id)
			if deviceBeforeSigning.LastSignature == deviceAfterSigning.LastSignature {
				t.Error("last signature value should be different after each sign operation")
			}
			if result == nil {
				t.Error("signature should be created, but it isn't")
			}
		}()
	}
	wg.Wait()
	deviceFromDb, _ := repository.GetById(id)
	if deviceFromDb.SignatureCounter != numOfSignatures {
		t.Errorf(
			"signature counter incorrect, got %d, expected %d",
			deviceFromDb.SignatureCounter,
			int64(numOfSignatures),
		)
	}
}

func testSignatureVerification(t *testing.T, algorithm string, temperedData string) bool {
	service := NewSignatureDeviceService(repository, locker)
	id := uuid.NewString()
	label := "Device"
	service.CreateSignatureDevice(id, algorithm, label)
	data := "message to be signed"
	signature, err := service.SignTransaction(id, data)
	if err != nil {
		t.Fatal("error occurred, test failed")
	}
	verified, _ := service.Verify(id, signature.Signature, signature.SignedData+temperedData)
	return verified.Status
}
