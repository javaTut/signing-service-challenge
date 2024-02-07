package repositories

import (
	"signing-service-challenge/domain"
	"signing-service-challenge/persistence"
)

type SignatureDeviceRepository interface {
	Save(domain.SignatureDevice) error
	GetById(string) (*domain.SignatureDevice, error)
	GetAll() ([]domain.SignatureDevice, error)
	// not in the requirements, but for testing purposes
	DeleteById(string) error
	DeleteAll() error
	Count() int
}

type SignatureDeviceInMemoryRepository struct {
	db persistence.InMemoryDB
}

func NewSignatureDeviceInMemoryRepository(db *persistence.InMemoryDB) *SignatureDeviceInMemoryRepository {
	return &SignatureDeviceInMemoryRepository{
		db: *db,
	}
}

func (r SignatureDeviceInMemoryRepository) Save(device domain.SignatureDevice) error {
	r.db.DevicesLock.Lock()
	defer r.db.DevicesLock.Unlock()
	r.db.Devices[device.Id] = device
	return nil
}

func (r SignatureDeviceInMemoryRepository) GetById(id string) (*domain.SignatureDevice, error) {
	// although it should not be its responsibility,
	// for the sake of simplicity, part of the locking logic is implemented here.
	r.db.DevicesLock.Lock()
	defer r.db.DevicesLock.Unlock()
	device, ok := r.db.Devices[id]
	if !ok {
		return nil, domain.ErrDeviceNotFound
	}
	return &device, nil
}

func (r SignatureDeviceInMemoryRepository) GetAll() ([]domain.SignatureDevice, error) {
	r.db.DevicesLock.RLock()
	defer r.db.DevicesLock.RUnlock()
	devices := []domain.SignatureDevice{}
	for _, value := range r.db.Devices {
		devices = append(devices, value)
	}
	return devices, nil
}

func (r SignatureDeviceInMemoryRepository) DeleteById(id string) error {
	r.db.DevicesLock.Lock()
	defer r.db.DevicesLock.Unlock()
	delete(r.db.Devices, id)
	return nil
}

func (r *SignatureDeviceInMemoryRepository) DeleteAll() error {
	r.db.DevicesLock.Lock()
	defer r.db.DevicesLock.Unlock()
	r.db.Devices = make(map[string]domain.SignatureDevice)
	return nil
}

func (r SignatureDeviceInMemoryRepository) Count() int {
	r.db.DevicesLock.RLock()
	defer r.db.DevicesLock.RUnlock()
	return len(r.db.Devices)
}
