package ipc

import (
	"sync"
)

type IPLock interface {
	sync.Locker
	RLock()
	RUnlock()
	Close()
}

type IPLockMode int8

const (
	SemLockMode IPLockMode = 0
	FLockMode   IPLockMode = 1
)

func NewIPLock(path string, mode IPLockMode, ftokId ...uint64) (IPLock, error) {
	if mode > 1 {
		mode = 0
	}
	switch mode {
	case FLockMode:
		flock, err := NewFlock(path)
		if err != nil {
			return nil, err
		}
		return flock.RWMutex(), nil
	case SemLockMode:
		var id uint64
		if len(ftokId) > 0 {
			id = ftokId[0]
		}
		key, err := Ftok(path, id)
		if err != nil {
			return nil, err
		}
		return NewSemLock(key)
	}
	return nil, nil
}
