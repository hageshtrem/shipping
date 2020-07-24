package infrastructure

import (
	"sync"
	
	"booking/domain"

	"github.com/pborman/uuid"
)

type cargoRepository struct{
	mtx sync.RWMutex
	cargos map[domain.TrackingID]*domain.Cargo
}

func (r *cargoRepository) Store(c *domain.Cargo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.cargos[c.TrackingID] = c
	return nil
}

func (r *cargoRepository) Find(id domain.TrackingID) (*domain.Cargo, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.cargos[id]; ok {
		return val, nil
	}
	return nil, domain.ErrUnknownCargo
}

func (r *cargoRepository) FindAll() []*domain.Cargo {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	c := make([]*domain.Cargo, 0, len(r.cargos))
	for _, val := range r.cargos {
		c = append(c, val)
	}
	return c
}

func (cargoRepository) NextTrackingID() domain.TrackingID {
	return domain.TrackingID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoRepository() domain.CargoRepository {
	return &cargoRepository{
		cargos: make(map[shipping.TrackingID]*shipping.Cargo),
	}
}

