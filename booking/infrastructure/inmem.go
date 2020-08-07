package infrastructure

import (
	"strings"
	"sync"

	"booking/domain"

	"github.com/pborman/uuid"
)

type cargoRepository struct {
	mtx    sync.RWMutex
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

func (*cargoRepository) NextTrackingID() domain.TrackingID {
	return domain.TrackingID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoRepository() domain.CargoRepository {
	return &cargoRepository{
		cargos: make(map[domain.TrackingID]*domain.Cargo),
	}
}

type locationRepository struct {
	locations map[domain.UNLocode]*domain.Location
}

func (r *locationRepository) Find(locode domain.UNLocode) (*domain.Location, error) {
	if l, ok := r.locations[locode]; ok {
		return l, nil
	}
	return nil, domain.ErrUnknownLocation
}

func (r *locationRepository) FindAll() []*domain.Location {
	l := make([]*domain.Location, 0, len(r.locations))
	for _, val := range r.locations {
		l = append(l, val)
	}
	return l
}

// NewLocationRepository returns a new instance of a in-memory location repository.
func NewLocationRepository() domain.LocationRepository {
	r := &locationRepository{
		locations: make(map[domain.UNLocode]*domain.Location),
	}

	r.locations[domain.SESTO] = domain.Stockholm
	r.locations[domain.AUMEL] = domain.Melbourne
	r.locations[domain.CNHKG] = domain.Hongkong
	r.locations[domain.USNYC] = domain.NewYork
	r.locations[domain.USCHI] = domain.Chicago
	r.locations[domain.JNTKO] = domain.Tokyo
	r.locations[domain.DEHAM] = domain.Hamburg
	r.locations[domain.NLRTM] = domain.Rotterdam
	r.locations[domain.FIHEL] = domain.Helsinki

	return r
}
