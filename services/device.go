package services

import (
	"encoding/base64"
	"fmt"
	"signing-service-challenge/crypto"
	"signing-service-challenge/domain"
	"signing-service-challenge/dto"
	"signing-service-challenge/lockers"
	"signing-service-challenge/repositories"
)

type SignatureDeviceService struct {
	repository repositories.SignatureDeviceRepository
	locker     lockers.DeviceLocker
}

func NewSignatureDeviceService(repository repositories.SignatureDeviceRepository, locker lockers.DeviceLocker) *SignatureDeviceService {
	return &SignatureDeviceService{
		repository: repository,
		locker:     locker,
	}
}

func (sd *SignatureDeviceService) CreateSignatureDevice(id, algorithm, label string) (*dto.CreateSignatureDeviceResponse, error) {
	device := domain.NewSignatureDeviceWithoutKeys(id, algorithm, label)
	kpHandler, err := crypto.GenerateKeyPairHandler(algorithm)
	if err != nil {
		return nil, err
	}
	privateKey, publicKey, err := kpHandler.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	kpHandler.AttachKeyPair(device, privateKey, publicKey)
	sd.repository.Save(*device)
	response := dto.ConvertSignatureDeviceToCreateResponse(*device)
	return &response, nil
}

func (sd *SignatureDeviceService) SignTransaction(deviceId string, data string) (*dto.SignatureResponse, error) {
	sd.locker.Lock(deviceId)
	defer sd.locker.Unlock(deviceId)
	// time.Sleep(1 * time.Millisecond)
	device, err := sd.repository.GetById(deviceId)
	if err != nil {
		return nil, domain.ErrDeviceNotFound
	}
	keyHandler, err := crypto.GenerateKeyPairHandler(device.Algorithm)
	if err != nil {
		return nil, err
	}
	primaryKey, err := keyHandler.Unmarshal(device.PrivateKey)
	if err != nil {
		return nil, err
	}
	signer, err := crypto.GenerateSigner(primaryKey)
	if err != nil {
		return nil, err
	}
	//<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
	securedDataToBeSigned := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, device.LastSignature)
	sign, err := signer.Sign([]byte(securedDataToBeSigned))
	if err != nil {
		return nil, err
	}
	signatureEncoded := base64.StdEncoding.EncodeToString(sign)
	device.LastSignature = signatureEncoded
	device.SignatureCounter = device.SignatureCounter + 1
	err = sd.repository.Save(*device)
	if err != nil {
		return nil, err
	}
	return &dto.SignatureResponse{
		Signature:  signatureEncoded,
		SignedData: securedDataToBeSigned,
	}, nil
}

func (sd *SignatureDeviceService) Verify(deviceId, signature, data string) (dto.VerificationResponse, error) {
	sd.locker.Lock(deviceId)
	defer sd.locker.Unlock(deviceId)
	device, err := sd.repository.GetById(deviceId)
	if err != nil {
		return dto.ConvertVerificationToResponse(false), err
	}
	keyHandler, err := crypto.GenerateKeyPairHandler(device.Algorithm)
	primaryKey, err := keyHandler.Unmarshal(device.PrivateKey)
	if err != nil {
		return dto.ConvertVerificationToResponse(false), err
	}
	signer, err := crypto.GenerateSigner(primaryKey)
	if err != nil {
		return dto.ConvertVerificationToResponse(false), err
	}
	sgn, err := base64.StdEncoding.DecodeString(signature)
	verified, err := signer.Verify([]byte(data), sgn)
	if err != nil {
		return dto.ConvertVerificationToResponse(false), err
	}
	return dto.ConvertVerificationToResponse(verified), nil
}

func (sd *SignatureDeviceService) GetById(deviceId string) (*dto.SignatureDeviceResponse, error) {
	device, err := sd.repository.GetById(deviceId)
	if err != nil {
		return nil, err
	}
	response := dto.ConvertSignatureDeviceToResponse(*device)
	return &response, nil
}

func (sd *SignatureDeviceService) GetAll() ([]dto.SignatureDeviceResponse, error) {
	devices, err := sd.repository.GetAll()
	if err != nil {
		return []dto.SignatureDeviceResponse{}, err
	}
	response := []dto.SignatureDeviceResponse{}
	for _, device := range devices {
		response = append(response, dto.ConvertSignatureDeviceToResponse(device))
	}
	return response, nil
}
