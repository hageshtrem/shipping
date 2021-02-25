package infrastructure

import (
	"sort"
	"strings"
	"sync"

	"booking/domain"

	"github.com/pborman/uuid"
)

type orderedCargo struct {
	id int
	*domain.Cargo
}

type orderedCargoSlice []orderedCargo

func (o orderedCargoSlice) Len() int {
	return len(o)
}

func (o orderedCargoSlice) Less(i, j int) bool {
	return o[i].id < o[j].id
}

func (o orderedCargoSlice) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type cargoRepository struct {
	mtx    sync.RWMutex
	lastID int
	cargos map[domain.TrackingID]orderedCargo
}

func (r *cargoRepository) Store(c *domain.Cargo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if oc, ok := r.cargos[c.TrackingID]; ok {
		r.cargos[c.TrackingID] = orderedCargo{oc.id, c}
	} else {
		r.cargos[c.TrackingID] = orderedCargo{r.lastID + 1, c}
		r.lastID++
	}
	return nil
}

func (r *cargoRepository) Find(id domain.TrackingID) (*domain.Cargo, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.cargos[id]; ok {
		return val.Cargo, nil
	}
	return nil, domain.ErrUnknownCargo
}

func (r *cargoRepository) FindAll() []*domain.Cargo {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	o := make(orderedCargoSlice, 0, len(r.cargos))
	for _, val := range r.cargos {
		o = append(o, val)
	}
	sort.Sort(o)
	c := make([]*domain.Cargo, 0, len(o))
	for _, val := range o {
		c = append(c, val.Cargo)
	}
	return c
}

func (*cargoRepository) NextTrackingID() domain.TrackingID {
	return domain.TrackingID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoRepository() domain.CargoRepository {
	return &cargoRepository{
		lastID: 0,
		cargos: make(map[domain.TrackingID]orderedCargo),
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
	r.locations[domain.SEGOT] = domain.Goteborg
	r.locations[domain.AUMEL] = domain.Melbourne
	r.locations[domain.CNHKG] = domain.Hongkong
	r.locations[domain.CNSHA] = domain.Shanghai
	r.locations[domain.CNHGH] = domain.Hangzhou
	r.locations[domain.USNYC] = domain.NewYork
	r.locations[domain.USCHI] = domain.Chicago
	r.locations[domain.USDAL] = domain.Dallas
	r.locations[domain.JNTKO] = domain.Tokyo
	r.locations[domain.DEHAM] = domain.Hamburg
	r.locations[domain.NLRTM] = domain.Rotterdam
	r.locations[domain.FIHEL] = domain.Helsinki

	return r
}
