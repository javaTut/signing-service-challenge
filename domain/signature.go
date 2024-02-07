package domain

type Signature struct {
	Id        string
	Signature string
	Data      string
	SignedBy  string
}

func NewSignature(id, signature, data, deviceId string) *Signature {
	return &Signature{
		Id:        id,
		Signature: signature,
		Data:      data,
		SignedBy:  deviceId,
	}
}
