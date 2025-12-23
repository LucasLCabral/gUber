package repository

import (
	"context"
	"fmt"
	"guber/services/trip-service/internal/domain"
)

type inMemoryRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemRepository() *inMemoryRepository {
	return &inMemoryRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inMemoryRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	fare, exists := r.rideFares[id]
	if !exists {
		return nil, fmt.Errorf("fare does not exist with the ID: %s", id)
	}
	return fare, nil
}

func (r *inMemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inMemoryRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.rideFares[f.ID.Hex()] = f
	
	return nil
}
