package dto

type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm" validate:"required"`
	Label     string `json:"label" validate:"required"`
}

type CreateSignatureDeviceResponse struct {
	Id               string `json:"id"`
	Algorithm        string `json:"algorithm"`
	Label            string `json:"label"`
	PrivateKey       string `json:"private_key"`
	PublicKey        string `json:"public_key"`
	SignatureCounter int    `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
}

type SignatureDeviceResponse struct {
	Id               string `json:"id"`
	Algorithm        string `json:"algorithm"`
	Label            string `json:"label"`
	PublicKey        string `json:"public_key"`
	SignatureCounter int    `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
}

type SignatureRequest struct {
	Id   string `json:"device_id" validate:"required"`
	Data string `json:"data" validate:"required"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type SignatureFullResponse struct {
	Id         string `json:"signature_id"`
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
	SignedBy   string `json:"signed_by"`
}

type VerificationResponse struct {
	Status bool `json:"status"`
}

type VerificationRequest struct {
	DeviceId  string `json:"device_id"`
	Signature string `json:"signature"`
	Data      string `json:"signed_data"`
}
