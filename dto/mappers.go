package dto

import "signing-service-challenge/domain"

func ConvertSignatureDeviceToCreateResponse(device domain.SignatureDevice) CreateSignatureDeviceResponse {
	// A private key is accessible only initially after the creation of a device
	return CreateSignatureDeviceResponse{
		Id:               device.Id,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		PrivateKey:       string(device.PrivateKey),
		PublicKey:        string(device.PublicKey),
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
	}
}

func ConvertSignatureDeviceToResponse(device domain.SignatureDevice) SignatureDeviceResponse {
	return SignatureDeviceResponse{
		Id:               device.Id,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		PublicKey:        string(device.PublicKey),
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
	}
}

func ConvertSignatureToResponse(signature domain.Signature) SignatureFullResponse {
	return SignatureFullResponse{
		Id:         signature.Id,
		Signature:  signature.Signature,
		SignedData: signature.Data,
		SignedBy:   signature.SignedBy,
	}
}

func ConvertVerificationToResponse(verification bool) VerificationResponse {
	return VerificationResponse{
		Status: verification,
	}
}
