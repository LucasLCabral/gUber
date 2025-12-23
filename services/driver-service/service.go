package main

import (
	pb "guber/shared/proto/driver"
	"guber/shared/util"
	math "math/rand/v2"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type Service struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// todo: route
}

func NewServive() *Service {
	return &Service{
		drivers: make([]*driverInMap, 0),
	}
}


func (s *Service) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	randomAvatar := util.GetRandomAvatar(randomIndex)
	randomPlate := GenerateRandomPlate()

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		GeoHash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		Id:             driverId,
		PackageSlug:    packageSlug,
		CarPlate:       randomPlate,
		ProfilePicture: randomAvatar,
	}

	s.drivers = append(s.drivers, &driverInMap{
		Driver: driver,
	})

	return driver, nil
}

func (s *Service) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, d := range s.drivers {
		if d.Driver.Id == driverId {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
			break
		}
	}
}
