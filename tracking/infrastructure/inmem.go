package infrastructure

import (
	"sync"
	"tracking/application"
)

type cargoRepository struct {
	mtx    sync.RWMutex
	cargos map[string]*application.Cargo
}

func (r *cargoRepository) Store(c *application.Cargo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.cargos[c.TrackingID] = c
	return nil
}

func (r *cargoRepository) Find(id string) (*application.Cargo, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.cargos[id]; ok {
		return val, nil
	}
	return nil, application.ErrUnknownCargo
}

func (r *cargoRepository) FindAll() []*application.Cargo {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	c := make([]*application.Cargo, 0, len(r.cargos))
	for _, val := range r.cargos {
		c = append(c, val)
	}
	return c
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoViewModelRepository() application.CargoViewModelRepository {
	return &cargoRepository{
		cargos: make(map[string]*application.Cargo),
	}
}
