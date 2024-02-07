package lockers

import (
	"sync"
)

type DeviceLocker interface {
	Lock(id string)
	Unlock(id string) error
}

// ID locking is required based on the assumption that multiple clients
// can simultaneously access the same device. If we use global lock in this
// case it will block all devices.
// The implementation of locking is not necessary if a device can only be
// used/changed by a single client at a time.

// Device locker implementation based on the map data structure (Set of locked IDs protected by
// global mutex and sync.Cond). It locks devices already in use based on their IDs.
// Credits: https://medium.com/@petrlozhkin/kmutex-lock-mutex-by-unique-id-408467659c24
// Other solutions, such as one with multiple mutexes are also possible, but this one is
// more light weight and good enough for this use case.

type DeviceLockerWithGlobalMapProtection struct {
	cond   *sync.Cond
	locker sync.Locker
	ids    map[string]struct{}
}

func NewDeviceLockerWithGlobalMapProtection(l sync.Locker) *DeviceLockerWithGlobalMapProtection {
	return &DeviceLockerWithGlobalMapProtection{
		cond:   sync.NewCond(l),
		locker: l,
		ids:    make(map[string]struct{}),
	}
}

func (p *DeviceLockerWithGlobalMapProtection) isLocked(id string) bool {
	_, ok := p.ids[id]
	return ok
}

func (p *DeviceLockerWithGlobalMapProtection) Lock(id string) {
	p.locker.Lock()
	defer p.locker.Unlock()
	for p.isLocked(id) {
		// wait for unlock
		p.cond.Wait()
	}
	p.ids[id] = struct{}{}
	return
}

func (p *DeviceLockerWithGlobalMapProtection) Unlock(id string) error {
	p.locker.Lock()
	defer p.locker.Unlock()
	delete(p.ids, id)
	// wake other clients waiting on this device
	p.cond.Broadcast()
	return nil
}
