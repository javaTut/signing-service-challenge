package services

import (
	"signing-service-challenge/domain"
	"signing-service-challenge/dto"
	"signing-service-challenge/repositories"
)

type SignatureService struct {
	repository repositories.SignatureRepository
}

func NewSignatureService(repository repositories.SignatureRepository) *SignatureService {
	return &SignatureService{
		repository: repository,
	}
}

func (sd SignatureService) Save(signature domain.Signature) error {
	sd.repository.Save(signature)
	return nil
}

func (sd SignatureService) GetById(signatureId string) (*dto.SignatureFullResponse, error) {
	signature, err := sd.repository.GetById(signatureId)
	if err != nil {
		return nil, err
	}
	response := dto.ConvertSignatureToResponse(*signature)
	return &response, nil
}

func (sd SignatureService) GetAll() ([]dto.SignatureFullResponse, error) {
	signatures, err := sd.repository.GetAll()
	if err != nil {
		return []dto.SignatureFullResponse{}, err
	}
	response := []dto.SignatureFullResponse{}
	for _, signature := range signatures {
		response = append(response, dto.ConvertSignatureToResponse(signature))
	}
	return response, nil
}
