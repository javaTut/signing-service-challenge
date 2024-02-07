package dto

import "errors"

func ValidateCreateSignatureDeviceRequest(request CreateSignatureDeviceRequest) (bool, error) {
	if request.Algorithm == "" {
		return false, errors.New("algorithm field is required")
	}
	if request.Label == "" {
		return false, errors.New("label field is required")
	}
	return true, nil
}

func ValidateSignRequest(request SignatureRequest) (bool, error) {
	if request.Id == "" {
		return false, errors.New("device_id field is required")
	}
	if request.Data == "" {
		return false, errors.New("data field is required")
	}
	return true, nil
}

func ValidateVerifyRequest(request VerificationRequest) (bool, error) {
	if request.DeviceId == "" {
		return false, errors.New("device_id field is required")
	}
	if request.Data == "" {
		return false, errors.New("data field is required")
	}
	if request.Signature == "" {
		return false, errors.New("data field is required")
	}
	return true, nil
}
