package repositories

import (
	"signing-service-challenge/domain"
	"signing-service-challenge/persistence"
)

type SignatureRepository interface {
	Save(domain.Signature) error
	GetById(string) (*domain.Signature, error)
	GetAll() ([]domain.Signature, error)
	// not in the requirements, but for test cleanups
	DeleteById(string) error
	DeleteAll() error
	Count() int
}

type SignatureInMemoryRepository struct {
	db persistence.InMemoryDB
}

func NewSignatureInMemoryRepository(db *persistence.InMemoryDB) *SignatureInMemoryRepository {
	return &SignatureInMemoryRepository{
		db: *db,
	}
}

func (r SignatureInMemoryRepository) Save(signature domain.Signature) error {
	r.db.SignaturesLock.Lock()
	defer r.db.SignaturesLock.Unlock()
	r.db.Signatures[signature.Id] = signature
	return nil
}

func (r SignatureInMemoryRepository) GetById(id string) (*domain.Signature, error) {
	r.db.SignaturesLock.Lock()
	defer r.db.SignaturesLock.Unlock()
	signature, ok := r.db.Signatures[id]
	if !ok {
		return nil, domain.ErrSignatureNotFound
	}
	return &signature, nil
}

func (r SignatureInMemoryRepository) GetAll() ([]domain.Signature, error) {
	r.db.SignaturesLock.RLock()
	defer r.db.SignaturesLock.RUnlock()
	signatures := []domain.Signature{}
	for _, value := range r.db.Signatures {
		signatures = append(signatures, value)
	}
	return signatures, nil
}

func (r SignatureInMemoryRepository) DeleteById(id string) error {
	r.db.SignaturesLock.Lock()
	defer r.db.SignaturesLock.Unlock()
	delete(r.db.Signatures, id)
	return nil
}

func (r *SignatureInMemoryRepository) DeleteAll() error {
	r.db.SignaturesLock.Lock()
	defer r.db.SignaturesLock.Unlock()
	r.db.Signatures = make(map[string]domain.Signature)
	return nil
}

func (r SignatureInMemoryRepository) Count() int {
	r.db.SignaturesLock.RLock()
	defer r.db.SignaturesLock.RUnlock()
	return len(r.db.Signatures)
}
