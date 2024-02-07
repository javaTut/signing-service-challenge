package main

import (
	"log"
	"sync"

	"signing-service-challenge/api"
	"signing-service-challenge/lockers"
	"signing-service-challenge/persistence"
	"signing-service-challenge/repositories"
	"signing-service-challenge/services"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	var db = persistence.NewInMemoryDB()

	// repositories
	deviceRepo := repositories.NewSignatureDeviceInMemoryRepository(db)
	signatureRepo := repositories.NewSignatureInMemoryRepository(db)

	// services
	var mutex sync.Mutex
	locker := lockers.NewDeviceLockerWithGlobalMapProtection(&mutex)
	deviceSvc := services.NewSignatureDeviceService(deviceRepo, locker)
	signatureSvc := services.NewSignatureService(signatureRepo)

	server := api.NewServer(ListenAddress, *deviceSvc, *signatureSvc)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
